package getload

import (
	"github.com/shirou/gopsutil/mem"
)

type MemoryObject struct {
	Buffered  uint64
	Total     uint64
	Free      uint64
	Used      uint64
	Swap      uint64
	Available uint64
}

func MemoryStat() MemoryObject {
	mbuf := make(chan uint64)
	mtotal := make(chan uint64)
	mfree := make(chan uint64)
	mused := make(chan uint64)
	mswap := make(chan uint64)
	mavailable := make(chan uint64)

	defer close(mbuf)
	defer close(mtotal)
	defer close(mfree)
	defer close(mused)
	defer close(mswap)
	defer close(mavailable)

	memory, err := mem.VirtualMemory()
	if err != nil {
		return MemoryObject{
			Buffered:  0,
			Total:     0,
			Free:      0,
			Used:      0,
			Swap:      0,
			Available: 0,
		}
	}

	go getMemoryBuffered(mbuf, memory)
	go getMemoryTotal(mtotal, memory)
	go getMemoryFree(mfree, memory)
	go getMemoryUsed(mused, memory)
	go getMemorySwap(mswap, memory)
	go getMemoryAvailable(mavailable, memory)

	return MemoryObject{
		Buffered:  <-mbuf,
		Total:     <-mtotal,
		Free:      <-mfree,
		Used:      <-mused,
		Swap:      <-mswap,
		Available: <-mavailable,
	}
}

func getMemoryBuffered(response chan<- uint64, memory *mem.VirtualMemoryStat) {
	response <- memory.Buffers
}

func getMemoryTotal(response chan<- uint64, memory *mem.VirtualMemoryStat) {
	response <- memory.Total
}

func getMemoryFree(response chan<- uint64, memory *mem.VirtualMemoryStat) {
	response <- memory.Free
}

func getMemoryUsed(response chan<- uint64, memory *mem.VirtualMemoryStat) {
	response <- memory.Used
}

func getMemorySwap(response chan<- uint64, memory *mem.VirtualMemoryStat) {
	response <- memory.SwapTotal
}

func getMemoryAvailable(response chan<- uint64, memory *mem.VirtualMemoryStat) {
	response <- memory.Available
}
