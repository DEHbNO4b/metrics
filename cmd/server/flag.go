package main

import (
	"flag"
	"os"
)

var flagRunAddr string

func parseFlag() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "adress and port for running")
	flag.Parse()
	if ep := os.Getenv("ADRESS"); ep != "" {
		flagRunAddr = ep
	}

}
