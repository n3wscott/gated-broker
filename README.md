# LEDHouse Broker

[![Go Report Card](https://goreportcard.com/badge/github.com/n3wscott/ledhouse-broker)](https://goreportcard.com/report/github.com/n3wscott/ledhouse-broker)
[![Godoc documentation](https://img.shields.io/badge/godoc-documentation-blue.svg)](https://godoc.org/github.com/n3wscott/ledhouse-broker/pkg)


This is a broker implementation for communication with the
[LEDHouse](http://github.com/n3wscott/ledhouse). It will serve as a demo for my
talk at Kubecon EU 2018 entitled
[Kubernetes as an Abstraction Layer for a Connected Home](http://sched.co/DqwC)

This leverages the [osb-broker-lib](github.com/pmorie/osb-broker-lib) to
bootstrap a broker. This broker is intended to be run locally on a laptop for
the talk. It is a representation of an on-prem service management.

The [k8s-broker-proxy](https://github.com/n3wscott/k8s-broker-proxy) is intended
to be paired with this to bridge it to the cloud/k8s environment.

TODO:
[] supply pub/sub connection info.
[] return pub/sub details in registry binding.
[] add flow to readme.
