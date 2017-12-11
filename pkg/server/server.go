package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/n3wscott/gated-broker/pkg/client"
	"github.com/n3wscott/gated-broker/pkg/controller"
	"github.com/n3wscott/osb-framework-go/pkg/apis/broker/v2"
	osb "github.com/n3wscott/osb-framework-go/pkg/server"
)

type server struct {
	Router     *mux.Router
	Controller v2.BrokerController
	hub        *client.Hub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func CreateServer() *server {

	hub := client.NewHub()

	go hub.Run()

	broker := controller.NewBrokerController(hub)

	osbServer := osb.CreateServer(broker)

	s := server{
		Router:     osbServer.Router,
		Controller: osbServer.Controller,
		hub:        hub,
	}

	s.Router.HandleFunc("/", s.GetHome).Methods("GET")
	s.Router.HandleFunc("/ws", s.GetWS)

	return &s
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *client.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client.NewClient(conn, hub)
}

func (s *server) GetWS(w http.ResponseWriter, r *http.Request) {
	serveWs(s.hub, w, r)
}
