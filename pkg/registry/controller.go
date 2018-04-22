package registry

import (
	"github.com/golang/glog"
	"github.com/n3wscott/ledhouse-broker/pkg/lightboard"
)

func NewControllerInstance(port string, lights map[Location]map[Kind]int) *ControllerInstance {
	c := ControllerInstance{}

	c.populateLightInstancesForLEDHouse(10)

	lightBoard, err := lightboard.NewLightBoard(port, 10)
	if err != nil {
		glog.Fatal("Failed to connect to light board on port ", port)
	}

	c.LightBoard = lightBoard
	return &c
}
