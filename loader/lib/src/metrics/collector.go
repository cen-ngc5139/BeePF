package metrics

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/cilium/ebpf"
)

// Collector 用于收集 eBPF 程序运行时统计信息
type Collector interface {
	// Start 启动收集器
	Start() error

	// Stop 停止收集器
	Stop() error

	// GetPrograms 获取当前所有 BPF 程序信息
	GetPrograms() ([]*meta.ProgramStats, error)

	// GetProgramStats 获取指定程序的统计信息
	GetProgramStats(id uint32) (*meta.MetricsStats, error)

	SetAttachedPros(map[uint32]*ebpf.Program) error

	GetAttachedPros() map[uint32]*ebpf.Program

	Export() error
}

// collector 实现了 Collector 接口
type StatsCollector struct {
	// 互斥锁保护并发访问
	mu sync.RWMutex

	// 所有 BPF 程序的映射
	programs map[uint32]*meta.ProgramStats

	// 当前已经挂载的 program 实例
	attachedPros map[uint32]*ebpf.Program

	// 统计数据缓存
	stats map[uint32]*meta.MetricsStats

	// 采集间隔
	interval time.Duration

	// 是否正在运行
	running bool

	// 停止信号通道
	stopCh chan struct{}

	// BPF stats 的关闭器
	closer io.Closer

	// 导出器
	exporterHandler meta.MetricsHandler

	// 日志记录器
	logger *zap.Logger
}

// NewCollector 创建一个新的收集器实例
func NewStatsCollector(interval time.Duration, exporterHandler meta.MetricsHandler, logger *zap.Logger) (Collector, error) {
	// 检查内核版本并启用 BPF stats
	closer, err := EnableBPFStats()
	if err != nil {
		return nil, fmt.Errorf("enable bpf stats failed: %v", err)
	}

	c := &StatsCollector{
		programs:        make(map[uint32]*meta.ProgramStats),
		stats:           make(map[uint32]*meta.MetricsStats),
		interval:        interval,
		stopCh:          make(chan struct{}),
		closer:          closer,
		exporterHandler: exporterHandler,
		logger:          logger,
	}

	return c, nil
}

func (c *StatsCollector) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return fmt.Errorf("collector is already running")
	}

	// 启动后台采集协程
	go c.Collect()

	c.running = true
	return nil
}

func (c *StatsCollector) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return nil
	}

	close(c.stopCh)
	c.running = false

	// 关闭 BPF stats
	if c.closer != nil {
		c.closer.Close()
	}

	return nil
}

func (c *StatsCollector) GetPrograms() ([]*meta.ProgramStats, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	programs := make([]*meta.ProgramStats, 0, len(c.programs))
	for _, p := range c.programs {
		programs = append(programs, p.Clone())
	}
	return programs, nil
}

func (c *StatsCollector) GetProgramStats(id uint32) (*meta.MetricsStats, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats, ok := c.stats[id]
	if !ok {
		return nil, fmt.Errorf("program %d not found", id)
	}

	return stats, nil
}

// collect 执行实际的数据采集
func (c *StatsCollector) Collect() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			if err := c.UpdateStats(); err != nil {
				// TODO: 错误处理
				continue
			}
		}
	}
}

// updateStats 更新所有程序的统计信息
func (c *StatsCollector) UpdateStats() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 更新程序信息和统计数据
	for _, prog := range c.attachedPros {
		info, err := prog.Info()
		if err != nil {
			return err
		}

		id, ok := info.ID()
		if !ok {
			return fmt.Errorf("fail to get prog %s id", info.Name)
		}

		// 更新或创建程序信息
		program, ok := c.programs[uint32(id)]
		if !ok {
			program = meta.NewProgram(prog)
			c.programs[uint32(id)] = program
		}
		program.Update(prog)

		// 更新统计信息
		stats, ok := c.stats[uint32(id)]
		if !ok {
			stats = meta.NewStats()
			c.stats[uint32(id)] = stats
		}
		stats.Update(program)
	}

	return nil
}

func (c *StatsCollector) SetAttachedPros(attached map[uint32]*ebpf.Program) error {
	if attached == nil {
		return errors.Errorf("failed to set attached pros, attached is nil")
	}

	c.attachedPros = attached
	return nil
}

func (c *StatsCollector) Export() error {
	if c.exporterHandler == nil {
		return errors.Errorf("failed to export stats, exporter handler is nil")
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-c.stopCh:
				return
			case <-ticker.C:
				programs, err := c.GetPrograms()
				if err != nil {
					c.logger.Error("获取 stats 信息失败", zap.Error(err))
				}

				for _, program := range programs {
					stats, err := c.GetProgramStats(program.ID)
					if err != nil {
						c.logger.Error("获取 stats 信息失败", zap.Error(err))
					}

					if err := c.exporterHandler.Handle(stats); err != nil {
						c.logger.Error("导出 stats 信息失败", zap.Error(err))
					}
				}
			}
		}
	}()

	return nil
}

func (c *StatsCollector) GetAttachedPros() map[uint32]*ebpf.Program {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.attachedPros
}
