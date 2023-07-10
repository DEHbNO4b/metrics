package main

import (
	"runtime"
	"sync"

	"github.com/DEHbNO4b/metrics/internal/agent"
)

<<<<<<< HEAD
var repInterval int = 10
var pInterval int = 2
var adress string = "localhost:8080"

=======
>>>>>>> 1892b541e7db1e891886b57fcbec69b4310e7abd
func main() {
	var m runtime.MemStats

	go agent.ReadRuntimeMetrics(&m, pInterval)
	go agent.PullMetrics(&m, repInterval, adress)
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
