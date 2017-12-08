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
		Controller: &controller.BrokerController{},
		hub:        newHub(),
	}

	go s.hub.run()

	s.Router.HandleFunc("/", s.GetHome).Methods("GET")
	s.Router.HandleFunc("/ws", s.GetWS)

	s.Router.HandleFunc("/v2/catalog", s.GetCatalog).Methods("GET")

	s.Router.HandleFunc("/v2/service_instances/{instance_id}", s.ProvisionServiceInstance).Methods("PUT")
	s.Router.HandleFunc("/v2/service_instances/{instance_id}", s.UpdateServiceInstance).Methods("PATCH")
	s.Router.HandleFunc("/v2/service_instances/{instance_id}", s.DeleteServiceInstance).Methods("DELETE")

	s.Router.HandleFunc("/v2/service_instances/{instance_id}last_operation", s.PollServiceInstance).Methods("GET")

	s.Router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", s.CreateServiceBinding).Methods("PUT")
	s.Router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", s.DeleteServiceBinding).Methods("DELETE")

	return &s
}
