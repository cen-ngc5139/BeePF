package models

import "time"

type MetricPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Value       float64   `json:"value"`
	ProgramName string    `json:"program_name"`
}

type TaskMetrics struct {
	AvgRunTimeNS      []MetricPoint `json:"avg_run_time_ns"`
	CPUUsage          []MetricPoint `json:"cpu_usage"`
	EventsPerSecond   []MetricPoint `json:"events_per_second"`
	PeriodNS          []MetricPoint `json:"period_ns"`
	TotalAvgRunTimeNS []MetricPoint `json:"total_avg_run_time_ns"`
}
