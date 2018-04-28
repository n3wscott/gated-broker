package broker

import (
	"flag"
)

type Options struct {
	SerialPort string

	// For Pub/Sub
	ProjectID    string
	Subscription string

	// Alt Pub/Sub config from a binding file.
	Binding string
}

// AddFlags is a hook called to initialize the CLI flags for broker options.
// It is called after the flags are added for the skeleton and before flag
// parse is called.
func AddFlags(o *Options) {
	flag.StringVar(&o.SerialPort, "serial", "", "The serial port to use to connect to the ledhouse.")

	flag.StringVar(&o.ProjectID, "projectId", "", "GCP projectId, for Pub/Sub")
	flag.StringVar(&o.Subscription, "subscription", "", "Pub/Sub subscription the Light Registry will use")

	flag.StringVar(&o.Binding, "binding", "", "Pub/Sub binding to use from Service Catalog")
}
