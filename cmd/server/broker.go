package server

import (
	"net/http"

	"encoding/json"
)

func (s *server) GetCatalog(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	catalog, _ := s.Controller.GetCatalog()
	json.NewEncoder(w).Encode(catalog)

	s.hub.broadcast <- []byte(string(catalog))
}
