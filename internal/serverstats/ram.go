package serverstats

import (
	"runtime"

	"github.com/pbnjay/memory"
)

type RAMStats struct {
	Sys   uint64 `json:"sys"`
	Used  uint64 `json:"used"`
	Total uint64 `json:"total"`
}

type RAM struct{}

func (RAM) GetStats() RAMStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return RAMStats{
		Sys:   bToMb(m.Sys),
		Used:  bToMb(m.TotalAlloc),
		Total: memory.TotalMemory(),
	}
}

func bToMb(b uint64) uint64 {
	return b / MB
}
