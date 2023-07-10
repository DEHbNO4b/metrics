package main

import (
	"runtime"
	"sync"

	"github.com/DEHbNO4b/metrics/internal/agent"
)

var repInterval int = 10
var pInterval int = 2
var adress string = "localhost:8080"

func main() {
	var m runtime.MemStats

	go agent.ReadRuntimeMetrics(&m, pInterval)
	go agent.PullMetrics(&m, repInterval, adress)
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
