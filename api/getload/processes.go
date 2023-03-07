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

func ProcStat() ProcObject {

	prunning := make(chan bool)
	ppid := make(chan int32)
	pname := make(chan string)

	defer close(prunning)
	defer close(ppid)
	defer close(pname)

	var procList []Process

	processes, err := process.Processes()
	if err != nil {
		log.Print(err.Error())
		return ProcObject{Processes: procList}
	}
	var wg sync.WaitGroup
	for _, proc := range processes[1:] {
		wg.Add(1)
		go func(proc *process.Process) {
			defer wg.Done()
			go getRunningProc(prunning, proc)
			go getProcPID(ppid, proc)
			go getProcName(pname, proc)
			procList = append(procList, Process{
				PID:     <-ppid,
				Running: <-prunning,
				Name:    <-pname,
			})
		}(proc)
		wg.Wait()
	}

	return ProcObject{Processes: procList}
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
