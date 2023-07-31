package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	endpoint       string
	reportInterval int
	pollInterval   int
	storeInterval  int
	filepath       string
	restore        bool
)

func parseFlag() {
	flag.StringVar(&endpoint, "a", "localhost:8080", "endpoint adress")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.IntVar(&storeInterval, "i", 300, "store interval")
	flag.StringVar(&filepath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(&restore, "r", true, "restore")
	flag.Parse()

	if ep := os.Getenv("ADDRESS"); ep != "" {
		endpoint = ep
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
