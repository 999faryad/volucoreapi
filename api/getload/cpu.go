// Package getload provides functions for obtaining CPU load and speed statistics.
package getload

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// CPUObject represents CPU statistics including model, speed, and load.
type CPUObject struct {
	Model   string  // CPU model name.
	Speed   float64 // CPU speed in MHz.
	Load    string  // CPU load as a percentage string.
	Cores   int32
	Threads int32
}

// CPUStat returns a CPUObject with CPU model name, speed, and load.
func CPUStat() CPUObject {
	// Use channels to concurrently retrieve CPU model name, speed, and load.
	cpumod := make(chan string)
	cpuload := make(chan string)
	cpuspeed := make(chan float64)
	cputhreads := make(chan int32)
	cpucores := make(chan int32)
	// defer closing all channels
	defer close(cpumod)
	defer close(cpuload)
	defer close(cpuspeed)
	defer close(cputhreads)
	defer close(cpucores)

	go getCPUModel(cpumod)
	go getCPULoad(cpuload)
	go getCPUSpeed(cpuspeed)
	go getCPUThreads(cputhreads)
	go getCPUCores(cpucores)

	return CPUObject{
		Model:   <-cpumod,
		Load:    <-cpuload,
		Speed:   <-cpuspeed,
		Cores:   <-cpucores,
		Threads: <-cputhreads,
	}

}

// getCPUModel retrieves the CPU model name.
func getCPUModel(response chan<- string) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Print(err.Error())
		response <- ""
		return
	}
	response <- strings.TrimSpace(cpuInfo[0].ModelName)
}

// getCPULoad retrieves the CPU load as a percentage string.
func getCPULoad(response chan<- string) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		response <- err.Error()
		fmt.Print(err.Error())
		return
	}
	response <- fmt.Sprintf("%.2f%%", percent[0])
}

// getCPUSpeed retrieves the CPU speed in MHz.
func getCPUSpeed(response chan<- float64) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Print(err.Error())
		response <- 0
		return
	}
	response <- cpuInfo[0].Mhz
}

func getCPUThreads(response chan<- int32) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Print(err.Error())
		response <- 0
		return
	}
	response <- cpuInfo[0].Cores
}
func getCPUCores(response chan<- int32) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Print(err.Error())
		response <- 0
		return
	}
	response <- cpuInfo[0].Cores / 2
}
