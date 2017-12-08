package v1

type Broker interface {
	CatalogController
	InstanceController
	BindingController
}
