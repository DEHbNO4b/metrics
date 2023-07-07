package main

import (
	"runtime"
	"sync"

	"github.com/DEHbNO4b/metrics/internal/agent"
)

type gauge struct {
	name string
	val  float64
}
type counter struct {
	name string
	val  int
}

func main() {
	var m runtime.MemStats

	go agent.ReadRuntimeMetrics(&m)
	go agent.PullMetrics(&m)
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
