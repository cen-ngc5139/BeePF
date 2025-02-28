package meta

import "time"

// MetricsStats 表示 BPF 程序的运行时统计信息
type MetricsStats struct {
	// CPU 使用率
	CPUTimePercent float64

	// 每秒事件数
	EventsPerSecond int64

	// 平均运行时间(ns)
	AvgRunTimeNS uint64

	// 总平均运行时间(ns)
	TotalAvgRunTimeNS uint64

	// 采样周期(ns)
	PeriodNS uint64

	// 最后更新时间
	LastUpdate time.Time
}

// NewStats 创建新的统计信息实例
func NewStats() *MetricsStats {
	return &MetricsStats{
		LastUpdate: time.Now(),
	}
}

// Update 根据程序信息更新统计数据
func (s *MetricsStats) Update(prog *ProgramStats) {
	now := time.Now()
	period := now.Sub(s.LastUpdate)
	s.PeriodNS = uint64(period.Nanoseconds())

	// 计算 CPU 使用率
	if s.PeriodNS > 0 {
		runtimeDelta := prog.RunTimeNS - prog.PrevRunTime
		s.CPUTimePercent = float64(runtimeDelta) / float64(s.PeriodNS) * 100
	}

	// 计算每秒事件数
	if s.PeriodNS > 0 {
		countDelta := prog.RunCount - prog.PrevCount
		s.EventsPerSecond = int64(float64(countDelta) / period.Seconds())
	}

	// 计算平均运行时间
	if countDelta := prog.RunCount - prog.PrevCount; countDelta > 0 {
		runtimeDelta := prog.RunTimeNS - prog.PrevRunTime
		s.AvgRunTimeNS = runtimeDelta / countDelta
	}

	// 计算总平均运行时间
	if prog.RunCount > 0 {
		s.TotalAvgRunTimeNS = prog.RunTimeNS / prog.RunCount
	}

	s.LastUpdate = now
}

// Clone 创建统计信息的副本
func (s *MetricsStats) Clone() *MetricsStats {
	clone := *s
	return &clone
}
