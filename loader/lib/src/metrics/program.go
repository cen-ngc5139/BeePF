package metrics

import (
	"time"

	"github.com/cilium/ebpf"
)

// Program 表示一个 BPF 程序
type Program struct {
	// 程序 ID
	ID uint32

	// 程序类型
	Type string

	// 程序名称
	Name string

	// 运行时统计
	RunTimeNS   uint64
	PrevRunTime uint64

	// 运行次数统计
	RunCount  uint64
	PrevCount uint64

	// 最后更新时间
	LastUpdate time.Time

	// 关联的进程信息
	Processes []Process
}

// Process 表示使用 BPF 程序的进程
type Process struct {
	PID  int32
	Comm string
}

// NewProgram 从 ebpf.Program 创建 Program
func NewProgram(prog *ebpf.Program) *Program {
	info, err := prog.Info()
	if err != nil {
		return nil
	}

	id, ok := info.ID()
	if !ok {
		return nil
	}

	return &Program{
		ID:         uint32(id),
		Type:       info.Type.String(),
		Name:       info.Name,
		LastUpdate: time.Now(),
	}
}

// Update 更新程序信息
func (p *Program) Update(prog *ebpf.Program) {
	// 保存上一次的统计数据
	p.PrevRunTime = p.RunTimeNS
	p.PrevCount = p.RunCount

	// 更新新的统计数据
	info, err := prog.Info()
	if err != nil {
		return
	}

	runtime, ok := info.Runtime()
	if !ok {
		return
	}

	runCount, ok := info.RunCount()
	if !ok {
		return
	}
	p.RunTimeNS = uint64(runtime)
	p.RunCount = runCount
	p.LastUpdate = time.Now()
}

// Clone 创建程序信息的副本
func (p *Program) Clone() *Program {
	clone := *p

	// 深拷贝进程信息
	if p.Processes != nil {
		clone.Processes = make([]Process, len(p.Processes))
		copy(clone.Processes, p.Processes)
	}

	return &clone
}
