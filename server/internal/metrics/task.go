package metrics

import (
	"fmt"
	"os"
	"sync"

	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type TaskMetrics struct {
	TaskStore *sync.Map
	TaskStats *TaskStatsMetrics
}

func NewTaskMetrics(store *sync.Map) *TaskMetrics {
	return &TaskMetrics{
		TaskStore: store,
		TaskStats: NewTaskStatsMetrics(),
	}
}

const (
	// 任务 CPU 使用率
	TaskCPUUsage = "beepf_task_cpu_usage"
	// 任务每秒事件数
	TaskEventsPerSecond = "beepf_task_events_per_second"
	// 任务平均运行时间(ns)
	TaskAvgRunTimeNS = "beepf_task_avg_run_time_ns"
	// 任务总平均运行时间(ns)
	TaskTotalAvgRunTimeNS = "beepf_task_total_avg_run_time_ns"
	// 任务采样周期(ns)
	TaskPeriodNS = "beepf_task_period_ns"
)

var (
	TaskMetricsLabels = []string{"task_id", "component_id", "program_id", "node_name"}
)

type TaskStatsMetrics struct {
	TaskCPUUsage          *prometheus.GaugeVec
	TaskEventsPerSecond   *prometheus.GaugeVec
	TaskAvgRunTimeNS      *prometheus.GaugeVec
	TaskTotalAvgRunTimeNS *prometheus.GaugeVec
	TaskPeriodNS          *prometheus.GaugeVec
}

func createGaugeVec(name, help string, labels []string) *prometheus.GaugeVec {
	return promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
}

func NewTaskStatsMetrics() *TaskStatsMetrics {
	return &TaskStatsMetrics{
		TaskCPUUsage:          createGaugeVec(TaskCPUUsage, "ebpf program task cpu usage", TaskMetricsLabels),
		TaskEventsPerSecond:   createGaugeVec(TaskEventsPerSecond, "ebpf program task events per second", TaskMetricsLabels),
		TaskAvgRunTimeNS:      createGaugeVec(TaskAvgRunTimeNS, "ebpf program task avg run time ns", TaskMetricsLabels),
		TaskTotalAvgRunTimeNS: createGaugeVec(TaskTotalAvgRunTimeNS, "ebpf program task total avg run time ns", TaskMetricsLabels),
		TaskPeriodNS:          createGaugeVec(TaskPeriodNS, "ebpf program task period ns", TaskMetricsLabels),
	}
}

func (m *TaskMetrics) ResetMetrics() {
	m.TaskStats.TaskCPUUsage.Reset()
	m.TaskStats.TaskEventsPerSecond.Reset()
	m.TaskStats.TaskAvgRunTimeNS.Reset()
	m.TaskStats.TaskTotalAvgRunTimeNS.Reset()
	m.TaskStats.TaskPeriodNS.Reset()
}

func (m *TaskMetrics) UpdateMetricsFromCache(nodeName string) {
	m.ResetMetrics()
	m.TaskStore.Range(func(key, value interface{}) bool {
		task := value.(*models.RunningTask)
		taskID := fmt.Sprintf("%d", task.Task.ID)
		for _, v := range task.Task.ProgStatus {
			// todo 此处通过 prog attach id 无法找到对应的 prog stats
			programStats, err := task.BPFLoader.StatsCollector.GetProgramStats(v.AttachID)
			if err != nil {
				fmt.Errorf("get program stats failed", zap.Error(err))
				return true
			}

			componentID := fmt.Sprintf("%d", task.Task.ComponentID)
			programID := fmt.Sprintf("%d", v.ID)
			m.TaskStats.TaskCPUUsage.WithLabelValues(taskID, componentID, programID, nodeName).Set(programStats.CPUTimePercent)
			m.TaskStats.TaskEventsPerSecond.WithLabelValues(taskID, componentID, programID, nodeName).Set(float64(programStats.EventsPerSecond))
			m.TaskStats.TaskAvgRunTimeNS.WithLabelValues(taskID, componentID, programID, nodeName).Set(float64(programStats.AvgRunTimeNS))
			m.TaskStats.TaskTotalAvgRunTimeNS.WithLabelValues(taskID, componentID, programID, nodeName).Set(float64(programStats.TotalAvgRunTimeNS))
			m.TaskStats.TaskPeriodNS.WithLabelValues(taskID, componentID, programID, nodeName).Set(float64(programStats.PeriodNS))
		}

		return true
	})
}

func (m *TaskMetrics) Handler() gin.HandlerFunc {
	h := promhttp.Handler()

	nodeName, err := os.Hostname()
	if err != nil {
		nodeName = "default_node"
	}

	return func(c *gin.Context) {
		m.UpdateMetricsFromCache(nodeName)
		h.ServeHTTP(c.Writer, c.Request)
	}
}
