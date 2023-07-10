package main

import "flag"

var flagRunAddr string

func parseFlag() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "adress and port for running")
	flag.Parse()
}
