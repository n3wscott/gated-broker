package controller

import (
	"github.com/n3wscott/gated-broker/pkg/client"
)

type BrokerController struct {
	hub *client.Hub
}

func NewBrokerController(hub *client.Hub) *BrokerController {
	c := BrokerController{hub: hub}
	return &c
}
