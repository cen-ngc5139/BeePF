package export

import (
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"go.uber.org/zap"
)

type MyCustomHandler struct {
	Logger *zap.Logger
}

// 实现 EventHandler 接口
func (h *MyCustomHandler) HandleEvent(ctx *meta.UserContext, data *meta.ReceivedEventData) error {
	switch data.Type {
	case meta.TypeJsonText:
		h.Logger.Info("received json data",
			zap.String("data", data.JsonText))
	case meta.TypePlainText:
		h.Logger.Info("received plain text",
			zap.String("data", data.Text))
	}
	return nil
}
