package LightRegistry

import (
	"fmt"

	"github.com/pborman/uuid"
)

func (c *ControllerInstance) AssignCredentials(osbInstanceId OsbId, osbBindingId OsbId) (*LightBinding, error) {
	// find the light id from the osb instance id
	lightId := c.OsbInstanceIdToId[osbInstanceId]
	if lightId == "" {
		return nil, fmt.Errorf("error: no light registered for OsbId[%s]", osbInstanceId)
	}

	// find the instance from the light id
	instance := c.IdToInstance[lightId]
	if instance == nil {
		return nil, fmt.Errorf("error: no light registered for OsbId[%s] and registery is in a bad state", osbInstanceId)
	}

	// make sure the binding id is not tied to a light
	if c.OsbBindingIdToId[osbBindingId] != "" {
		return nil, fmt.Errorf("error: bindingId[%s] in use", osbBindingId)
	}

	// make sure the instance does not have a conflicting binding
	for _, binding := range instance.Bindings {
		if binding.OsbBindingId == osbBindingId {
			return nil, fmt.Errorf("error: bindingId[%s] in use and we are in a bad state", osbBindingId)
		}
	}

	// ok to make a binding for the given instance
	binding := LightBinding{
		OsbBindingId: osbBindingId,
		Id:           lightId,
		Secret:       Secret(uuid.NewUUID().String()),
	}
	// TODO: confirm that the secret is unique?

	c.OsbBindingIdToId[osbBindingId] = lightId
	c.SecretToId[binding.Secret] = lightId
	instance.Bindings = append(instance.Bindings, binding)

	return &binding, nil
}

func (c *ControllerInstance) RemoveCredentials() (*string, error) {
	return nil, nil
}
