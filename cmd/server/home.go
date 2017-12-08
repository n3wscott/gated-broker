package server

import (
	"net/http"

	"github.com/golang/glog"
)

func (s *server) GetHome(w http.ResponseWriter, r *http.Request) {
	glog.Info(r.URL.Path)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "cmd/server/static/home.html")
}
