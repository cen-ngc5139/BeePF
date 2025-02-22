package metrics

import (
	"io"
	"os"

	"github.com/cilium/ebpf"
)

const (
	bpfStatsEnabled = "/proc/sys/kernel/bpf_stats_enabled"
	bpfStatsRunTime = 0 // BPF_STATS_RUN_TIME from linux/bpf.h
)

// EnableBPFStats 启用或禁用 BPF stats
// 在 Linux 内核版本 >= 5.8 上，使用 syscall 启用
// 在 Linux 内核版本 < 5.8 上，使用 procfs 启用
// 同时，只会对当前进程启动的 bpf 程序生效
func EnableBPFStats() (io.Closer, error) {
	// 尝试通过 syscall 启用 (内核版本 >= 5.8)
	fd, err := ebpf.EnableStats(bpfStatsRunTime)
	if err == nil {
		return fd, nil
	}

	// 回退到通过 procfs 启用
	return nil, os.WriteFile(bpfStatsEnabled, []byte("1"), 0644)
}
