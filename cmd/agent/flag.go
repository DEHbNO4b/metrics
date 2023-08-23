package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	endpoint       string
	key            string
	reportInterval int
	pollInterval   int
)

func parseFlag() {
	flag.StringVar(&endpoint, "a", "localhost:8080", "endpoint adress")
	flag.StringVar(&key, "k", "", "hash key")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.Parse()
	if ep := os.Getenv("ADDRESS"); ep != "" {
		endpoint = ep
	}
	if k := os.Getenv("KEY"); k != "" {
		key = k
	}
	if ri := os.Getenv("REPORT_INTERVAL"); ri != "" {
		rInt, err := strconv.Atoi(ri)
		if err != nil {
			reportInterval = rInt
		}

	}
	if pi := os.Getenv("POLL_INTERVAL"); pi != "" {
		pInt, err := strconv.Atoi(pi)
		if err != nil {
			pollInterval = pInt
		}
	}
}
