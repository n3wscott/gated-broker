package LightRegistry

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

	c.Register("aabbcc", "Bedroom", "Red")
	_, err := c.Register("aabbcc", "Bedroom", "Red")
	if err == nil {
		t.Errorf("Expected second call to register the same OSB id to fail.")
	}

	glog.Infof("After registering:\n%s", c)

	binding, err := c.AssignCredentials("aabbcc", "binding-aabbcc")

	glog.Infof("Got back this binding: %s", binding)

	c.SetLightIntensity(binding.Secret, .5)

	glog.Infof("After light toggle:\n%s", c)

	glog.Flush()
}
