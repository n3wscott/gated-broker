package broker

import (
	"flag"
)

type Options struct {
	SerialPort string
}

// AddFlags is a hook called to initialize the CLI flags for broker options.
// It is called after the flags are added for the skeleton and before flag
// parse is called.
func AddFlags(o *Options) {
	flag.StringVar(&o.SerialPort, "serial", "", "The serial port to use to connect to the ledhouse.")
}
