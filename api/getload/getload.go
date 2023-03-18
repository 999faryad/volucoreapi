package getload

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RequestData struct {
	Show []string `json:"show"`
}

type Load struct {
	CPU       *CPUObject    `json:"CPU,omitempty"`
	Disk      *DiskObject   `json:"Disk,omitempty"`
	Memory    *MemoryObject `json:"Memory,omitempty"`
	Processes *ProcObject   `json:"Processes,omitempty"`
}

func GetLoad(writer http.ResponseWriter, request *http.Request) {
	var payload RequestData
	load := &Load{}

	reqBody, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
	}

	json.Unmarshal(reqBody, &payload)

	if len(payload.Show) == 0 {
		load = &Load{
			CPU:       CPUStat(),
			Memory:    MemoryStat(),
			Disk:      DiskStat(),
			Processes: ProcStat(),
		}
	}

	for _, data := range payload.Show {
		switch data {
		case "cpu":
			load.CPU = CPUStat()
		case "disk":
			load.Disk = DiskStat()
		case "memory":
			load.Memory = MemoryStat()
		case "processes":
			load.Processes = ProcStat()
		default:
			break
		}
	}

	fmt.Fprintln(writer, load.out())

}

func (load Load) out() string {
	response, err := json.Marshal(load)
	if err != nil {
		log.Println("Failed to Marshal Request")
	}
	return string(response)
}

func CloseChannels(chans ...interface{}) {
	for _, ch := range chans {
		if c, ok := ch.(chan interface{}); ok {
			close(c)
		}
	}
}
