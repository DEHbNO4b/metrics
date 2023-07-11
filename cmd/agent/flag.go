package main

import (
	"flag"
	"os"
	"strconv"
)

var endpoint string
var reportInterval int
var pollInterval int

func parseFlag() {
	flag.StringVar(&endpoint, "a", "localhost:8080", "endpoint adress")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.Parse()
	if ep := os.Getenv("ADRESS"); ep != "" {
		endpoint = ep
	}
	if ri := os.Getenv("REPORT_INTERVAL"); ri != "" {
		rInt, err := strconv.Atoi(ri)
		if err != nil {
			reportInterval = rInt
		}

	}
	if pi := os.Getenv("ADRESS"); pi != "" {
		pInt, err := strconv.Atoi(pi)
		if err != nil {
			pollInterval = pInt
		}
	}
}
