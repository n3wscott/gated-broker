package server

import (
	"github.com/gorilla/mux"
	"github.com/n3wscott/gated-broker/pkg/apis/broker/v1"
	"github.com/n3wscott/gated-broker/pkg/controller"
)

type server struct {
	Router     *mux.Router
	Controller v1.Broker
	hub        *Hub
}

func CreateServer() *server {

	s := server{
		Router:     mux.NewRouter(),
		Controller: &controller.Broker{},
		hub:        newHub(),
	}

	go s.hub.run()

	s.Router.HandleFunc("/", s.GetHome).Methods("GET")
	s.Router.HandleFunc("/ws", s.GetWS)

	s.Router.HandleFunc("/v2/catalog", s.GetCatalog).Methods("GET")

	return &s
}
