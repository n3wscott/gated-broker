package broker

import (
	"testing"

	"fmt"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"github.com/pmorie/osb-broker-lib/pkg/broker"
	"gopkg.in/yaml.v2"
)

func TestBusinessLogic_GetCatalog(t *testing.T) {
	response := &broker.CatalogResponse{}

	data := `
---
services:
- name: example-starter-pack-service
  id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a
  description: The example service from the osb starter pack!
`
	//
	//	`
	//	---
	//	services:
	//	- name: example-starter-pack-service
	//	  id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a
	//	  description: The example service from the osb starter pack!
	//	  bindable: true
	//	  plan_updateable: true
	//	  metadata:
	//		displayName: "Example starter-pack service"
	//		imageUrl: https://avatars2.githubusercontent.com/u/19862012?s=200&v=4
	//	  plans:
	//	  - name: default
	//		id: 86064792-7ea2-467b-af93-ac9694d96d5b
	//		description: The default plan for the starter pack example service
	//		free: true
	//		schemas:
	//		  service_instance:
	//			create:
	//			  "$schema": "http://json-schema.org/draft-04/schema"
	//			  "type": "object"
	//			  "title": "Parameters"
	//			  "properties":
	//			  - "name":
	//				  "title": "Some Name"
	//				  "type": "string"
	//				  "maxLength": 63
	//				  "default": "My Name"
	//			  - "color":
	//				  "title": "Color"
	//				  "type": "string"
	//				  "default": "Clear"
	//				  "enum":
	//				  - "Clear"
	//				  - "Beige"
	//				  - "Grey"
	//	`

	err := yaml.Unmarshal([]byte(data), &response)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Broker: %+v\n\n\n", response)
}

func TestBusinessLogic_GetCatalog2(t *testing.T) {
	response := &osb.CatalogResponse{}

	data := `
---
services:
- name: example-starter-pack-service
  id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a
  description: The example service from the osb starter pack!
`
	//
	//	`
	//	---
	//	services:
	//	- name: example-starter-pack-service
	//	  id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a
	//	  description: The example service from the osb starter pack!
	//	  bindable: true
	//	  plan_updateable: true
	//	  metadata:
	//		displayName: "Example starter-pack service"
	//		imageUrl: https://avatars2.githubusercontent.com/u/19862012?s=200&v=4
	//	  plans:
	//	  - name: default
	//		id: 86064792-7ea2-467b-af93-ac9694d96d5b
	//		description: The default plan for the starter pack example service
	//		free: true
	//		schemas:
	//		  service_instance:
	//			create:
	//			  "$schema": "http://json-schema.org/draft-04/schema"
	//			  "type": "object"
	//			  "title": "Parameters"
	//			  "properties":
	//			  - "name":
	//				  "title": "Some Name"
	//				  "type": "string"
	//				  "maxLength": 63
	//				  "default": "My Name"
	//			  - "color":
	//				  "title": "Color"
	//				  "type": "string"
	//				  "default": "Clear"
	//				  "enum":
	//				  - "Clear"
	//				  - "Beige"
	//				  - "Grey"
	//	`

	err := yaml.Unmarshal([]byte(data), &response)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("OSB: %+v\n\n\n", response)
}

type CR osb.CatalogResponse
type CatalogResponse struct {
	osb.CatalogResponse
}

func TestBusinessLogic_GetCatalog3(t *testing.T) {
	response := &CatalogResponse{}

	data := `
---
custom: this is a string
services:
- name: example-starter-pack-service
  id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a
  description: The example service from the osb starter pack!
`
	err := yaml.Unmarshal([]byte(data), &response.CatalogResponse)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Custom: %+v\n\n\n", response)
}
