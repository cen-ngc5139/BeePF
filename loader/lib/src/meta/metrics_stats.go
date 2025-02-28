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
