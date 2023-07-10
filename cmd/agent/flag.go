package main

import "flag"

var endpoint string
var reportInterval int
var pollInterval int

func parseFlag() {
	flag.StringVar(&endpoint, "a", "localhost:8080", "endpoint adress")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.Parse()
}
