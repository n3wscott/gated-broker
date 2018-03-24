package broker

import (
	"sync"

	"github.com/pmorie/osb-broker-lib/pkg/broker"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/n3wscott/gated-broker/pkg/registry"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"gopkg.in/yaml.v2"
)

// NewBusinessLogic is a hook that is called with the Options the program is run
// with. NewBusinessLogic is the place where you will initialize your
// BusinessLogic the parameters passed in.
func NewBusinessLogic(o Options) (*BusinessLogic, error) {
	// For example, if your BusinessLogic requires a parameter from the command
	// line, you would unpack it from the Options and set it on the
	// BusinessLogic here.

	// TODO: light registry creation needs to happen somewhere else
	lights := map[registry.Location]map[registry.Kind]int{
		"Bedroom": {
			"Red":   3,
			"Green": 2,
			"Blue":  1,
		},
		"Kitchen": {
			"Red":   1,
			"Green": 2,
			"Blue":  3,
		},
	}

	return &BusinessLogic{
		async:    o.Async,
		Registry: registry.NewControllerInstance(lights),
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
	Registry registry.Controller
}

func (b *BusinessLogic) AdditionalRouting(router *mux.Router) {
	// TODO: could pass in the router to the registry and it can do the assigning internally.
	router.HandleFunc("/graph", b.Registry.HandleGetGraph).Methods("GET")
}

var _ broker.Interface = &BusinessLogic{}

func (b *BusinessLogic) GetCatalog(c *broker.RequestContext) (*broker.CatalogResponse, error) {
	// Your catalog business logic goes here
	response := &broker.CatalogResponse{}

	data := `
---
services:
- name: example-starter-pack-service
  id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a
  description: The example service from the osb starter pack!
  bindable: true
  plan_updateable: true
  metadata:
    displayName: "Example starter-pack service"
    imageUrl: https://avatars2.githubusercontent.com/u/19862012?s=200&v=4
  plans:
  - name: default
    id: 86064792-7ea2-467b-af93-ac9694d96d5b
    description: The default plan for the starter pack example service
    free: true
    schemas:
      service_instance:
        create:
          "$schema": "http://json-schema.org/draft-04/schema"
          "type": "object"
          "title": "Parameters"
          "properties":
          - "name":
              "title": "Some Name"
              "type": "string"
              "maxLength": 63
              "default": "My Name"
          - "color":
              "title": "Color"
              "type": "string"
              "default": "Clear"
              "enum":
              - "Clear"
              - "Beige"
              - "Grey"
`

	err := yaml.Unmarshal([]byte(data), &response.CatalogResponse)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (b *BusinessLogic) Provision(request *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
	// Your provision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := broker.ProvisionResponse{}

	//exampleInstance := &exampleInstance{ID: request.InstanceID, Params: request.Parameters}
	//b.instances[request.InstanceID] = exampleInstance

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
	// Your deprovision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := broker.DeprovisionResponse{}

	//delete(b.instances, request.InstanceID)

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

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

	//instance, ok := b.instances[request.InstanceID]
	//if !ok {
	//	return nil, osb.HTTPStatusCodeError{
	//		StatusCode: http.StatusNotFound,
	//	}
	//}

	response := broker.BindResponse{
	//Credentials: instance.Params,
	}
	if request.AcceptsIncomplete {
		response.Async = b.async
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

// example types

// exampleInstance is intended as an example of a type that holds information about a service instance
type exampleInstance struct {
	ID     string
	Params map[string]interface{}
}
