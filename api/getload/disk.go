package getload

import (
	"github.com/shirou/gopsutil/disk"
	"log"
	"strings"
	"sync"
)

type Disks struct {
	Device string
	IOPS   uint64
	Total  uint64
	Free   uint64
	Used   uint64
}

type DiskObject struct {
	Disk []Disks
}

func DiskStat() DiskObject {
	partitions, _ := disk.Partitions(false)
	var disks []Disks

	ddevice := make(chan string)
	diops := make(chan uint64)
	dtotal := make(chan uint64)
	dfree := make(chan uint64)
	dused := make(chan uint64)

	defer close(ddevice)
	defer close(diops)
	defer close(dtotal)
	defer close(dfree)
	defer close(dused)

	var wg sync.WaitGroup
	for _, partition := range partitions {
		if !strings.HasPrefix(partition.Device, "/dev/sd") {
			continue
		}

		wg.Add(1)
		go func(partition disk.PartitionStat) {
			defer wg.Done()
			go getDiskDevice(ddevice, partition)
			go getDiskIOPS(diops, partition)
			go getDiskTotal(dtotal, partition)
			go getDiskFree(dfree, partition)
			go getDiskUsed(dused, partition)
			disks = append(disks, Disks{
				Device: <-ddevice,
				IOPS:   <-diops,
				Total:  <-dtotal,
				Free:   <-dfree,
				Used:   <-dused,
			})
		}(partition)

		wg.Wait()
	}

	return DiskObject{Disk: disks}
}

func getDiskDevice(response chan<- string, part disk.PartitionStat) {
	response <- part.Device
}

func getDiskIOPS(response chan<- uint64, part disk.PartitionStat) {
	counters, err := disk.IOCounters()
	if err != nil {
		log.Print(err.Error())
		response <- 0
		return
	}

	response <- counters[part.Device].IopsInProgress
}

func getDiskTotal(response chan<- uint64, part disk.PartitionStat) {
	device, err := disk.Usage(part.Device)
	if err != nil {
		log.Print(err.Error())
		response <- 0
		return
	}
	response <- device.Total

}

func getDiskFree(response chan<- uint64, part disk.PartitionStat) {
	device, err := disk.Usage(part.Device)
	if err != nil {
		log.Print(err.Error())
		response <- 0
		return
	}
	response <- device.Free
}

func getDiskUsed(response chan<- uint64, part disk.PartitionStat) {
	device, err := disk.Usage(part.Device)
	if err != nil {
		log.Print(err.Error())
		response <- 0
		return
	}

	response <- device.Used
}
