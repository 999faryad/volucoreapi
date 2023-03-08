package getload

import (
	"github.com/shirou/gopsutil/process"
	"log"
	"sync"
)

type Process struct {
	PID     int32
	Running bool
	Name    string
}

type ProcObject struct {
	Processes []Process
}

func ProcStat() *ProcObject {

	pRunning := make(chan bool)
	pPID := make(chan int32)
	pName := make(chan string)

	defer CloseChannels(pRunning, pPID, pName)

	var processList []Process
	processes, err := process.Processes()
	if err != nil {
		log.Print(err.Error())
		return &ProcObject{Processes: processList}
	}

	var wg sync.WaitGroup
	for _, proc := range processes[1:] {
		wg.Add(1)
		go func(proc *process.Process) {
			defer wg.Done()
			go getRunningProc(pRunning, proc)
			go getProcPID(pPID, proc)
			go getProcName(pName, proc)
			processList = append(processList, Process{
				PID:     <-pPID,
				Running: <-pRunning,
				Name:    <-pName,
			})
		}(proc)
		wg.Wait()
	}

	return &ProcObject{Processes: processList}
}

func getRunningProc(response chan<- bool, proc *process.Process) {
	running, err := proc.IsRunning()
	if err != nil {
		log.Print(err.Error())
		response <- false
		return
	}

	response <- running
}

func getProcPID(response chan<- int32, proc *process.Process) {
	response <- proc.Pid
}

func getProcName(response chan<- string, proc *process.Process) {
	name, err := proc.Name()
	if err != nil {
		log.Print(err.Error())
		response <- ""
		return
	}

	response <- name
}
