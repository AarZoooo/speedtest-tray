package speedtest_util

import (
	"runtime"

	"speedtest-tray/internal/config"
)

type ProcessStats struct {
	AllocMB      uint64
	SysMB        uint64
	NumGoroutine int
}

func GetProcessStats() ProcessStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return ProcessStats{
		AllocMB:      m.Alloc / 1024 / 1024,
		SysMB:        m.Sys / 1024 / 1024,
		NumGoroutine: runtime.NumGoroutine(),
	}
}

func (ps ProcessStats) LogAttr() []interface{} {
	return []interface{}{
		config.KeyAllocMB, ps.AllocMB,
		config.KeySysMB, ps.SysMB,
		config.KeyNumGoroutine, ps.NumGoroutine,
	}
}
