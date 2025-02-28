package meta

// EventHandler 定义事件处理接口
type EventHandler interface {
	HandleEvent(ctx *UserContext, data *ReceivedEventData) error
}

// MetricsHandler 用于处理 eBPF 程序的运行时统计信息
type MetricsHandler interface {
	// Handle 处理统计信息
	Handle(stats *MetricsStats) error
}
