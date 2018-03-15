package LightRegistry

type Controller interface {
	RegistrationController
	CredentialsController
	CatalogController
}

type RegistrationController interface {
	Register(osbInstanceId OsbId, location Location, kind Kind) (*LightInstance, error)
	Deregister(osbInstanceId OsbId) error
}

type CredentialsController interface {
	AssignCredentials(osbInstanceId OsbId, osbBindingId OsbId) (*LightBinding, error)
	RemoveCredentials() (*string, error)
}

type CatalogController interface {
	GetCatalog() (*string, error)
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
	Id       LightId
	Location Location
	Kind     Kind
	RGBLight
	WhiteLight
	TemperatureLight
}

type RGBLight struct {
	Intensity float32
	Red       float32
	Green     float32
	Blue      float32
}

type WhiteLight struct {
	Intensity float32
}

type TemperatureLight struct {
	Intensity   float32
	Temperature float32 // percentage cool to warm
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
	// Master list.
	IdToInstance map[LightId]LightInstance

	// Light location+Kind
	LocationKindToId map[Location]map[Kind]LightId

	// Helpful lookup lists.
	OsbInstanceIdToId map[OsbId]LightId
	SecretToId        map[Secret]LightId
	OsbBindingIdToId  map[OsbId]LightId
}

func NewControllerInstance() *ControllerInstance {
	c := ControllerInstance{}
	return &c
}
