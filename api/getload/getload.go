package getload

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Load struct {
	CPU       CPUObject
	Disk      DiskObject
	Memory    MemoryObject
	Processes ProcObject
}

func GetLoad(writer http.ResponseWriter, request *http.Request) {
	load := Load{CPU: CPUStat(), Disk: DiskStat(), Memory: MemoryStat(), Processes: ProcStat()}
	response, err := json.Marshal(load)
	if err != nil {
		fmt.Fprint(writer, err.Error())
	}
	fmt.Fprint(writer, string(response))
}
