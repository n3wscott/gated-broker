package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/n3wscott/gated-broker/cmd/server"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	s := server.CreateServer()

	glog.Infof("Starting Broker, %s", "http://localhost:12345")
	glog.Fatal(http.ListenAndServe(":12345", s.Router))
}
