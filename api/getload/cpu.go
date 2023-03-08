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
	Load    float64 // CPU load as a percentage string.
	Cores   int32
	Threads int32
}

// CPUStat returns a CPUObject with CPU model name, speed, and load.
func CPUStat() *CPUObject {

	// Use channels to concurrently retrieve CPU model name, speed, and load.
	cpuModel := make(chan string)
	cpuLoad := make(chan float64)
	cpuSpeed := make(chan float64)
	cpuThreads := make(chan int32)
	cpuCores := make(chan int32)

	// defer closing all channels
	defer CloseChannels(cpuModel, cpuLoad, cpuSpeed, cpuThreads, cpuCores)

	go getCPUModel(cpuModel)
	go getCPULoad(cpuLoad)
	go getCPUSpeed(cpuSpeed)
	go getCPUThreads(cpuThreads)
	go getCPUCores(cpuCores)

	return &CPUObject{
		Model:   <-cpuModel,
		Load:    <-cpuLoad,
		Speed:   <-cpuSpeed,
		Cores:   <-cpuCores,
		Threads: <-cpuThreads,
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
func getCPULoad(response chan<- float64) {
	percent, err := cpu.Percent(time.Millisecond, false)
	if err != nil {
		response <- 0
		fmt.Print(err.Error())
		return
	}
	response <- percent[0]
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
