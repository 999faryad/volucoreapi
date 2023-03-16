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

func MemoryStat() *MemoryObject {
	mBuffered := make(chan uint64)
	mTotal := make(chan uint64)
	mFree := make(chan uint64)
	mUsed := make(chan uint64)
	mSwap := make(chan uint64)
	mAvailable := make(chan uint64)

	defer CloseChannels(mBuffered, mTotal, mFree, mUsed, mSwap, mAvailable)

	memory, err := mem.VirtualMemory()
	if err != nil {
		return &MemoryObject{
			Buffered:  0,
			Total:     0,
			Free:      0,
			Used:      0,
			Swap:      0,
			Available: 0,
		}
	}

	go getMemoryBuffered(mBuffered, memory)
	go getMemoryTotal(mTotal, memory)
	go getMemoryFree(mFree, memory)
	go getMemoryUsed(mUsed, memory)
	go getMemorySwap(mSwap, memory)
	go getMemoryAvailable(mAvailable, memory)

	return &MemoryObject{
		Buffered:  <-mBuffered,
		Total:     <-mTotal,
		Free:      <-mFree,
		Used:      <-mUsed,
		Swap:      <-mSwap,
		Available: <-mAvailable,
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
