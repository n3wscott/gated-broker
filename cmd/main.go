package main

import (
	"context"
	"flag"

	"fmt"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"

	"github.com/golang/glog"
	prom "github.com/prometheus/client_golang/prometheus"

	//"github.com/n3wscott/gated-broker/pkg/server"

	"github.com/n3wscott/gated-broker/pkg/broker"
	"github.com/n3wscott/gated-broker/pkg/registry"
	"github.com/pmorie/osb-broker-lib/pkg/metrics"
	"github.com/pmorie/osb-broker-lib/pkg/rest"
	"github.com/pmorie/osb-broker-lib/pkg/server"
)

var addr = flag.String("addr", ":8080", "http service address")

//
//func main() {
//	flag.Parse()
//
//	s := server.CreateServer()
//
//	glog.Infof("Starting Broker, %s", "http://localhost:12345")
//	glog.Fatal(http.ListenAndServe(":12345", s.Router))
//
//}

var options struct {
	broker.Options

	Port    int
	TLSCert string
	TLSKey  string
}

func init() {
	flag.IntVar(&options.Port, "port", 3000, "use '--port' option to specify the port for broker to listen on")
	flag.StringVar(&options.TLSCert, "tlsCert", "", "base-64 encoded PEM block to use as the certificate for TLS. If '--tlsCert' is used, then '--tlsKey' must also be used. If '--tlsCert' is not used, then TLS will not be used.")
	flag.StringVar(&options.TLSKey, "tlsKey", "", "base-64 encoded PEM block to use as the private key matching the TLS certificate. If '--tlsKey' is used, then '--tlsCert' must also be used")
	broker.AddFlags(&options.Options)
	flag.Parse()
}

func main() {
	if err := run(); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		glog.Fatalln(err)
	}
}

func run() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go cancelOnInterrupt(ctx, cancelFunc)

	return runWithContext(ctx)
}

func runWithContext(ctx context.Context) error {
	if flag.Arg(0) == "version" {
		fmt.Printf("%s/%s\n", path.Base(os.Args[0]), "0.1.0")
		return nil
	}
	if (options.TLSCert != "" || options.TLSKey != "") &&
		(options.TLSCert == "" || options.TLSKey == "") {
		fmt.Println("To use TLS, both --tlsCert and --tlsKey must be used")
		return nil
	}

	addr := ":" + strconv.Itoa(options.Port)

	businessLogic, err := broker.NewBusinessLogic(options.Options)
	if err != nil {
		return err
	}

	//// Prom. metrics
	reg := prom.NewRegistry()
	osbMetrics := metrics.New()
	reg.MustRegister(osbMetrics)

	api, err := rest.NewAPISurface(businessLogic, osbMetrics)
	if err != nil {
		return err
	}

	s := server.New(api, reg)

	glog.Infof("Starting broker!")

	// TODO: light registry creation needs to happen somewhere else
	lights := map[registry.Location]map[registry.Kind]int{
		"Bedroom": {
			"Red":   3,
			"Green": 2,
			"Blue":  1,
		},
		"Kitchen": {
			"Red":   1,
			"Green": 2,
			"Blue":  3,
		},
	}

	c := registry.NewControllerInstance(lights)
	c.Register("aabbcc", "Bedroom", "Red")
	binding, _ := c.AssignCredentials("aabbcc", "binding-aabbcc")
	c.SetLightIntensity(binding.Secret, .5)

	// TODO: could pass in the router to the registry and it can do the assigning internally.
	s.Router.HandleFunc("/graph", c.HandleGetGraph).Methods("GET")

	if options.TLSCert == "" && options.TLSKey == "" {
		err = s.Run(ctx, addr)
	} else {
		err = s.RunTLS(ctx, addr, options.TLSCert, options.TLSKey)
	}
	return err
}

func cancelOnInterrupt(ctx context.Context, f context.CancelFunc) {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-term:
			glog.Infof("Received SIGTERM, exiting gracefully...")
			f()
			os.Exit(0)
		case <-ctx.Done():
			os.Exit(0)
		}
	}
}
