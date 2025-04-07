package topology

import (
	"errors"
	"fmt"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"github.com/cilium/ebpf"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
)

type NodeMetricsCollector struct {
	Collector metrics.Collector
}

func NewNodeMetricsCollector(interval time.Duration, logger *zap.Logger) (*NodeMetricsCollector, error) {
	collector, err := metrics.NewStatsCollector(interval, nil, logger)
	if err != nil {
		return nil, fmt.Errorf("fail to create metrics collector: %v", err)
	}

	progs, err := GetNodePrograms()
	if err != nil {
		return nil, fmt.Errorf("fail to get node programs: %v", err)
	}

	collector.SetAttachedPros(progs)

	return &NodeMetricsCollector{
		Collector: collector,
	}, nil
}

// GetAllProgramsMetrics 获取所有程序的运行时统计信息
func GetNodePrograms() (map[uint32]*ebpf.Program, error) {
	progs := make(map[uint32]*ebpf.Program)
	var nextID ebpf.ProgramID

	// 从 ID 0 开始遍历所有程序
	for {
		// 获取下一个程序 ID
		id, err := ebpf.ProgramGetNextID(nextID)
		if err != nil {
			// 如果是 ENOENT 错误，表示已经遍历完所有程序
			if errors.Is(err, unix.ENOENT) {
				break
			}
			return nil, fmt.Errorf("获取下一个程序 ID 失败: %w", err)
		}

		// 确保 ID 有变化，避免无限循环
		if id <= nextID {
			break
		}

		prog, err := ebpf.NewProgramFromID(id)
		if err != nil {
			return nil, fmt.Errorf("fail to get ebpf prog %d, err: %v", id, err)
		}
		progs[uint32(id)] = prog

		// 更新下一个 ID
		nextID = id
	}

	return progs, nil
}

func (c *NodeMetricsCollector) Start() error {
	if err := c.Collector.Start(); err != nil {
		return fmt.Errorf("fail to start metrics collector: %v", err)
	}

	return nil
}

func (c *NodeMetricsCollector) Stop() error {
	progs := c.Collector.GetAttachedPros()
	for _, v := range progs {
		v.Close()
	}

	return c.Collector.Stop()
}

func (c *NodeMetricsCollector) GetMetrics() (map[uint32]*meta.ProgMetricsStats, error) {
	metrics := make(map[uint32]*meta.ProgMetricsStats)
	programs, err := c.Collector.GetPrograms()
	if err != nil {
		return nil, fmt.Errorf("fail to get programs: %v", err)
	}

	for _, program := range programs {
		stats, err := c.Collector.GetProgramStats(program.ID)
		if err != nil {
			return nil, fmt.Errorf("fail to get program stats: %v", err)
		}

		metrics[program.ID] = &meta.ProgMetricsStats{
			ID:    program.ID,
			Type:  program.Type,
			Name:  program.Name,
			Stats: stats,
		}
	}
	return metrics, nil
}
