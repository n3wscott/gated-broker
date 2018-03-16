package LightRegistry

import (
	"testing"

	"github.com/golang/glog"
)

func TestNewControllerInstanceTest(t *testing.T) {

	lights := map[Location]map[Kind]int{
		"AAA": {
			"Red":   3,
			"Green": 2,
			"Blue":  4,
		},
		"BBB": {
			"Red":   4,
			"Green": 5,
			"Blue":  6,
		},
		"CCC": {
			"Red":   1,
			"Green": 2,
			"Blue":  3,
		},
	}

	c := NewControllerInstance(lights)

	glog.Info(c)
}
