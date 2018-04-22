package broker

import (
	"sync"

	"github.com/pmorie/osb-broker-lib/pkg/broker"

	"strings"

	"net/http"

	"fmt"

	"reflect"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/n3wscott/gated-broker/pkg/registry"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

// NewBusinessLogic is a hook that is called with the Options the program is run
// with. NewBusinessLogic is the place where you will initialize your
// BusinessLogic the parameters passed in.
func NewBusinessLogic(o Options) (*BusinessLogic, error) {
	// For example, if your BusinessLogic requires a parameter from the command
	// line, you would unpack it from the Options and set it on the
	// BusinessLogic here.

	// TODO: light registry creation needs to happen somewhere else
	lights := make(map[registry.Location]map[registry.Kind]int, 10)

	return &BusinessLogic{
		async:     o.Async,
		instances: make(map[string]*Instance, 10),
		Registry:  registry.NewControllerInstance(o.SerialPort, lights),
	}, nil
}

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type BusinessLogic struct {
	// Indicates if the broker should handle the requests asynchronously.
	async bool
	// Synchronize go routines.
	sync.RWMutex
	// The light registry
	Registry *registry.ControllerInstance
	// todo
	instances map[string]*Instance

	catalog *broker.CatalogResponse
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

	// Your catalog business logic goes here
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
	// Your provision business logic goes here

	// example implementation:
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

	location, kind := b.osbServicePlanToRegistryLocationKind(request.ServiceID, request.PlanID)
	light, err := b.Registry.Register(registry.OsbId(request.InstanceID), location, kind)

	if err != nil {
		return nil, err
	}

	b.instances[request.InstanceID] = instance

	dashboardURL := fmt.Sprintf("http:///%s/", string(light.Id))
	response.DashboardURL = &dashboardURL

	// when we support async:
	//if request.AcceptsIncomplete {
	//	response.Async = b.async
	//}

	return &response, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
	// Your deprovision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := broker.DeprovisionResponse{}

	// TODO

	delete(b.instances, request.InstanceID)

	//if request.AcceptsIncomplete {
	//	response.Async = b.async
	//}

	return &response, nil
}

func (b *BusinessLogic) LastOperation(request *osb.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error) {
	// Your last-operation business logic goes here

	return nil, nil
}

func (b *BusinessLogic) Bind(request *osb.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error) {
	// Your bind business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	_, ok := b.instances[request.InstanceID]
	if !ok {
		return nil, osb.HTTPStatusCodeError{
			StatusCode: http.StatusNotFound,
		}
	}

	lightBinding, err := b.Registry.AssignCredentials(registry.OsbId(request.InstanceID), registry.OsbId(request.BindingID))
	if err != nil {
		return nil, err
	}
	//if request.AcceptsIncomplete {
	//	response.Async = b.async
	//}

	url := fmt.Sprintf("http://localhost:3000/light/%s", lightBinding.Secret)

	response := broker.BindResponse{
		BindResponse: osb.BindResponse{
			Credentials: map[string]interface{}{"url": url},
		},
	}

	return &response, nil
}

func (b *BusinessLogic) Unbind(request *osb.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error) {
	// Your unbind business logic goes here
	return &broker.UnbindResponse{}, nil
}

func (b *BusinessLogic) Update(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
	// Your logic for updating a service goes here.
	response := broker.UpdateInstanceResponse{}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
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
