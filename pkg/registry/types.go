package registry

import (
	"net/http"

	"github.com/n3wscott/gated-broker/pkg/lightboard"
)

type Controller interface {
	RegistrationController
	CredentialsController
	CatalogController
	LightController
	HttpController
}

type RegistrationController interface {
	Register(osbInstanceId OsbId, location Location, kind Kind) (*LightInstance, error)
	Deregister(osbInstanceId OsbId) error
}

type CredentialsController interface {
	AssignCredentials(osbInstanceId OsbId, osbBindingId OsbId) (*LightBinding, error)
	RemoveCredentials(osbBindingId OsbId) error
}

type CatalogController interface {
	GetCatalog() (*string, error)
}

type LightController interface {
	SetLightIntensity(secret Secret, intensity float32) error
}

type HttpController interface {
	Graph() string
	HandleGetGraph(w http.ResponseWriter, r *http.Request)
}

// The internal id of the light.
type LightId string

// The container of all the light instance details.
type LightInstance struct {
	OsbInstanceId OsbId
	Id            LightId
	Bindings      []LightBinding
}

type Light struct {
	Id        LightId
	Index     int
	Location  Location
	Kind      Kind
	Intensity float32
}

// The container of all the light binding details.
type LightBinding struct {
	OsbBindingId OsbId
	Id           LightId
	Secret       Secret
}

type Location string // service class
type Kind string     // service plan
type OsbId string
type Secret string

type ControllerInstance struct {

	// Master list of lights.
	IdToLight map[LightId]*Light

	// Master list of instances.
	IdToInstance map[LightId]*LightInstance

	// Lights for a Location+Kind
	LocationKindToIds map[Location]map[Kind][]LightId

	// Helpful lookup lists.
	OsbInstanceIdToId map[OsbId]LightId
	SecretToId        map[Secret]LightId
	OsbBindingIdToId  map[OsbId]LightId

	LightBoard *lightboard.LightBoard
}
