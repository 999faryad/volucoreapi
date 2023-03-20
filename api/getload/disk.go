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

	diskDevice := make(chan string)
	diskIOPS := make(chan uint64)
	diskTotal := make(chan uint64)
	diskFree := make(chan uint64)
	diskUsed := make(chan uint64)

	defer CloseChannels(diskDevice, diskIOPS, diskTotal, diskFree, diskUsed)

	var wg sync.WaitGroup
	for _, partition := range partitions {
		if !strings.HasPrefix(partition.Device, "/dev/sd") && !strings.HasPrefix(partition.Device, "/dev/nvme") {
			continue
		}

		wg.Add(1)
		go func(partition disk.PartitionStat) {
			defer wg.Done()
			go getDiskDevice(diskDevice, partition)
			go getDiskIOPS(diskIOPS, partition)
			go getDiskTotal(diskTotal, partition)
			go getDiskFree(diskFree, partition)
			go getDiskUsed(diskUsed, partition)
			disks = append(disks, Disks{
				Device: <-diskDevice,
				IOPS:   <-diskIOPS,
				Total:  <-diskTotal,
				Free:   <-diskFree,
				Used:   <-diskUsed,
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
