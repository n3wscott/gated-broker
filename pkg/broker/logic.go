package broker

import (
	"sync"

	"github.com/pmorie/osb-broker-lib/pkg/broker"

	"strings"

	"net/http"

	"fmt"

	"reflect"

	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/n3wscott/ledhouse-broker/pkg/binding"
	"github.com/n3wscott/ledhouse-broker/pkg/registry"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

const LightRegistryServiceID = "light-registry"

func usage() {
	fmt.Println("usage:")
	flag.PrintDefaults()
	os.Exit(-1)
}

func NewBusinessLogic(o Options) (*BusinessLogic, error) {
	// TODO: light registry creation should happen somewhere else
	lights := make(map[registry.Location]map[registry.Kind]int, 10)

	if o.SerialPort == "" {
		fmt.Println("Error: Serial Port required.")
		usage()
	}

	if o.Binding != "" {
		projectId, _, subscriptionId, err := binding.PubSubBinding(o.Binding)
		if err != nil {
			return nil, err
		}
		o.ProjectID = projectId
		o.Subscription = subscriptionId
	}

	return &BusinessLogic{
		instances: make(map[string]*Instance, 10),
		Registry:  registry.NewControllerInstance(o.SerialPort, lights, o.ProjectID, o.Subscription),
	}, nil
}

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type BusinessLogic struct {
	// Synchronize go routines.
	sync.RWMutex
	// The light registry
	Registry *registry.ControllerInstance
	// todo
	instances map[string]*Instance

	catalog *broker.CatalogResponse
	Url     string
}

var _ broker.Interface = &BusinessLogic{}

func (b *BusinessLogic) AdditionalRouting(router *mux.Router) {
	// TODO: could pass in the router to the registry and it can do the assigning internally.
	router.HandleFunc("/graph", b.Registry.HandleGetGraph).Methods("GET")
	router.HandleFunc("/light/{secret}", b.Registry.HandleSetLight).Methods("PUT")
}

func (b *BusinessLogic) GetCatalog(c *broker.RequestContext) (*broker.CatalogResponse, error) {
	if b.catalog != nil {
		return b.catalog, nil
	}

	response := &broker.CatalogResponse{}

	for location, kinds := range b.Registry.LocationKindToIds {
		service := osb.Service{
			ID:          strings.ToLower("location-" + string(location)),
			Name:        string(location),
			Description: "A set of lights in " + string(location),
			Bindable:    true,
		}
		for kind, _ := range kinds {
			plan := osb.Plan{
				ID:          strings.ToLower("location-" + string(location) + "-kind-" + string(kind)),
				Name:        string(kind),
				Description: "Light type " + string(kind),
			}
			service.Plans = append(service.Plans, plan)
		}
		response.Services = append(response.Services, service)
	}

	// Add the light registry.
	{
		service := osb.Service{
			ID:          LightRegistryServiceID,
			Name:        "Light Registry",
			Description: "The light registry, the endpoint to use light bindings.",
			Bindable:    true,
		}

		plan := osb.Plan{
			ID:          "default",
			Name:        "default",
			Description: "Default light registry",
		}
		service.Plans = append(service.Plans, plan)
		response.Services = append(response.Services, service)
	}

	// save the catalog for later
	b.catalog = response
	return response, nil
}

func (b *BusinessLogic) osbServicePlanToRegistryLocationKind(serviceId, planId string) (registry.Location, registry.Kind) {
	if b.catalog == nil {
		b.GetCatalog(nil) // todo this could throw an error in theory
	}
	for _, s := range b.catalog.Services {
		if s.ID == serviceId {
			for _, p := range s.Plans {
				if p.ID == planId {
					return registry.Location(s.Name), registry.Kind(p.Name)
				}
			}
		}
	}
	return registry.Location(""), registry.Kind("")
}

func (b *BusinessLogic) Provision(request *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
	b.Lock()
	defer b.Unlock()

	response := broker.ProvisionResponse{}
	instance := &Instance{
		ID:        request.InstanceID,
		ServiceID: request.ServiceID,
		PlanID:    request.PlanID,
		Params:    request.Parameters,
	}

	if b.instances[request.InstanceID] != nil {
		i := b.instances[request.InstanceID]
		if i.Match(instance) {
			response.Exists = true
		} else {
			glog.Error("InstanceID in use")
			return nil, fmt.Errorf("InstanceID in use")
		}
	}
	b.instances[request.InstanceID] = instance
	if request.ServiceID == LightRegistryServiceID {
		dashboardURL := fmt.Sprintf("%s/light/", b.Url)
		response.DashboardURL = &dashboardURL
	} else { // assume it is a light
		location, kind := b.osbServicePlanToRegistryLocationKind(request.ServiceID, request.PlanID)
		light, err := b.Registry.Register(registry.OsbId(request.InstanceID), location, kind)

		if err != nil {
			return nil, err
		}
		dashboardURL := fmt.Sprintf("http:///%s/", string(light.Id))
		response.DashboardURL = &dashboardURL
	}
	return &response, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
	b.Lock()
	defer b.Unlock()

	response := broker.DeprovisionResponse{}

	// TODO(n3wscott): Implement this.

	delete(b.instances, request.InstanceID)

	return &response, nil
}

func (b *BusinessLogic) LastOperation(request *osb.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error) {
	return nil, nil
}

func (b *BusinessLogic) Bind(request *osb.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error) {
	b.Lock()
	defer b.Unlock()

	instance, ok := b.instances[request.InstanceID]
	if !ok {
		return nil, osb.HTTPStatusCodeError{
			StatusCode: http.StatusNotFound,
		}
	}

	response := broker.BindResponse{}
	if instance.ServiceID == LightRegistryServiceID {
		url := fmt.Sprintf("%s/light/{token}", b.Url)
		response.BindResponse = osb.BindResponse{
			Credentials: map[string]interface{}{
				"topic": b.Registry.Subscription.ID(), // I know these do not match but I am going to claim that that are same name.
				"url":   url,
			},
		}
	} else { // assume it is a light
		lightBinding, err := b.Registry.AssignCredentials(registry.OsbId(request.InstanceID), registry.OsbId(request.BindingID))
		if err != nil {
			return nil, err
		}
		url := fmt.Sprintf("%s/light/%s", b.Url, lightBinding.Secret)
		response.BindResponse = osb.BindResponse{
			Credentials: map[string]interface{}{
				"token": lightBinding.Secret,
				"url":   url,
			},
		}
	}
	return &response, nil
}

func (b *BusinessLogic) Unbind(request *osb.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error) {
	// TODO(n3wscott): Implement this.
	return &broker.UnbindResponse{}, nil
}

func (b *BusinessLogic) Update(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
	// Not supported.
	return &broker.UpdateInstanceResponse{}, nil
}

func (b *BusinessLogic) ValidateBrokerAPIVersion(version string) error {
	glog.Info("ValidateBrokerAPIVersion")
	return nil
}

type Instance struct {
	ID        string
	ServiceID string
	PlanID    string
	Params    map[string]interface{}
}

func (i *Instance) Match(other *Instance) bool {
	return reflect.DeepEqual(i, other)
}
