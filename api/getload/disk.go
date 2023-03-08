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

func DiskStat() *DiskObject {
	partitions, _ := disk.Partitions(false)
	var disks []Disks

	dDevice := make(chan string)
	dIOPS := make(chan uint64)
	dTotal := make(chan uint64)
	dFree := make(chan uint64)
	dUsed := make(chan uint64)

	defer CloseChannels(dDevice, dIOPS, dTotal, dFree, dUsed)

	var wg sync.WaitGroup
	for _, partition := range partitions {
		if !strings.HasPrefix(partition.Device, "/dev/sd") || !strings.HasPrefix(partition.Device, "/dev/nvme") {
			continue
		}

		wg.Add(1)
		go func(partition disk.PartitionStat) {
			defer wg.Done()
			go getDiskDevice(dDevice, partition)
			go getDiskIOPS(dIOPS, partition)
			go getDiskTotal(dTotal, partition)
			go getDiskFree(dFree, partition)
			go getDiskUsed(dUsed, partition)
			disks = append(disks, Disks{
				Device: <-dDevice,
				IOPS:   <-dIOPS,
				Total:  <-dTotal,
				Free:   <-dFree,
				Used:   <-dUsed,
			})
		}(partition)

		wg.Wait()
	}

	return &DiskObject{Disk: disks}
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
