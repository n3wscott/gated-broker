package server

import (
	"net/http"

	"encoding/json"
)

func (s *server) GetCatalog(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	catalog, _ := s.Controller.GetCatalog()
	json.NewEncoder(w).Encode(catalog)

	jsonCatalog, _ := json.Marshal(catalog)
	s.hub.broadcast <- []byte(jsonCatalog)
}

func (s *server) ProvisionServiceInstance(w http.ResponseWriter, req *http.Request) {

}

func (s *server) UpdateServiceInstance(w http.ResponseWriter, req *http.Request) {

}

func (s *server) DeleteServiceInstance(w http.ResponseWriter, req *http.Request) {

}

func (s *server) PollServiceInstance(w http.ResponseWriter, req *http.Request) {

}

func (s *server) CreateServiceBinding(w http.ResponseWriter, req *http.Request) {

}

func (s *server) DeleteServiceBinding(w http.ResponseWriter, req *http.Request) {

}
