package meta

// EventHandler 定义事件处理接口
type EventHandler interface {
	HandleEvent(ctx *UserContext, data *ReceivedEventData) error
}
