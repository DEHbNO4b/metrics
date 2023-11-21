package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/DEHbNO4b/metrics/internal/config"
)

var (
	endpoint       string
	key            string
	crypto         string
	reportInterval int
	pollInterval   int
)

func parseFlag() {
	flag.StringVar(&endpoint, "a", "localhost:8080", "endpoint adress")
	flag.StringVar(&key, "k", "", "hash key")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.StringVar(&crypto, "crypto-key", "", "crypto config file path")
	flag.Parse()
	if ep := os.Getenv("ADDRESS"); ep != "" {
		endpoint = ep
	}
	if ck := os.Getenv("CRYPTO_KEY"); ck != "" {
		crypto = ck
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
	if crypto != "" {
		_, err := config.DecPub(crypto)
		if err != nil {
			fmt.Println(err)
		}
	}
}
