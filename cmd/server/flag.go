package main

import "flag"

var flagRunAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "adress and port to run server")
	flag.Parse()
}
