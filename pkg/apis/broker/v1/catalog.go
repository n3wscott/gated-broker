package v1

type CatalogController interface {
	GetCatalog() (*Catalog, error)
}

type Catalog struct {
	Services []*Service `json:"services"`
}

type Service struct {
	Name            string        `json:"name"`
	ID              string        `json:"id"`
	Description     string        `json:"description"`
	Tags            []string      `json:"tags,omitempty"`
	Requires        []string      `json:"requires,omitempty"`
	Bindable        bool          `json:"bindable"`
	Metadata        interface{}   `json:"metadata,omitempty"`
	DashboardClient interface{}   `json:"dashboard_client"`
	PlanUpdateable  bool          `json:"plan_updateable,omitempty"`
	Plans           []ServicePlan `json:"plans"`
}

type ServicePlan struct {
	Name        string      `json:"name"`
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Metadata    interface{} `json:"metadata,omitempty"`
	Free        bool        `json:"free,omitempty"`
	Bindable    *bool       `json:"bindable,omitempty"`
	Schemas     *Schemas    `json:"schemas,omitempty"`
}

// Schemas represents a plan's schemas for service instance and binding create
// and update.
type Schemas struct {
	ServiceInstances *ServiceInstanceSchema `json:"service_instance,omitempty"`
	ServiceBindings  *ServiceBindingSchema  `json:"service_binding,omitempty"`
}

// ServiceInstanceSchema represents a plan's schemas for a create and update
// of a service instance.
type ServiceInstanceSchema struct {
	Create *InputParameters `json:"create,omitempty"`
	Update *InputParameters `json:"update,omitempty"`
}

// ServiceBindingSchema represents a plan's schemas for the parameters
// accepted for binding creation.
type ServiceBindingSchema struct {
	Create *InputParameters `json:"create,omitempty"`
}

// InputParameters represents a schema for input parameters for creation or
// update of an API resource.
type InputParameters struct {
	Parameters interface{} `json:"parameters,omitempty"`
}
