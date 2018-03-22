package registry

import (
	"testing"

	"github.com/golang/glog"
)

func TestNewControllerInstanceTest(t *testing.T) {

	lights := map[Location]map[Kind]int{
		"Bedroom": {
			"Red":   3,
			"Green": 2,
			"Blue":  4,
		},
		"LivingRoom": {
			"Red":   4,
			"Green": 5,
			"Blue":  6,
		},
		"Kitchen": {
			"Red":   1,
			"Green": 2,
			"Blue":  3,
		},
	}

	c := NewControllerInstance(lights)

	glog.Infof("Starting state:\n%s", c)

	instance, err := c.Register("aabbcc", "Bedroom", "Red")

	glog.Infof("Graph: \n\n%s\n\n", instance.Graph())

	_, err = c.Register("aabbcc", "Bedroom", "Red")
	if err == nil {
		t.Errorf("expected second call to register the same OSB id to fail")
	}

	glog.Infof("after registering:\n%s", c)

	binding, err := c.AssignCredentials("aabbcc", "binding-aabbcc")

	glog.Infof("got back this binding: %s", binding)

	c.SetLightIntensity(binding.Secret, .5)
	c.SetLightIntensity(binding.Secret, .7)

	glog.Infof("after light toggle:\n%s", c)

	if err := c.Deregister("aabbcc"); err == nil {
		t.Errorf("expected Deregister to fail, we still have bindings")
	}

	c.RemoveCredentials(binding.OsbBindingId)

	c.Deregister("aabbcc")

	glog.Infof("after deregister:\n%s", c)

	glog.Flush()
}
