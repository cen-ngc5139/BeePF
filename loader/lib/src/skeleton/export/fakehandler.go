package export

import "go.uber.org/zap"

type MyCustomHandler struct {
	Logger *zap.Logger
}

// 实现 EventHandler 接口
func (h *MyCustomHandler) HandleEvent(ctx *UserContext, data *ReceivedEventData) error {
	switch data.Type {
	case TypeJsonText:
		h.Logger.Info("received json data",
			zap.String("data", data.JsonText))
	case TypePlainText:
		h.Logger.Info("received plain text",
			zap.String("data", data.Text))
	}
	return nil
}
