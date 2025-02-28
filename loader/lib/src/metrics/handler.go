package metrics

import (
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"go.uber.org/zap"
)

type DefaultHandler struct {
	Logger *zap.Logger
}

func NewDefaultHandler(logger *zap.Logger) *DefaultHandler {
	return &DefaultHandler{
		Logger: logger,
	}
}

func (h *DefaultHandler) Handle(stats *meta.MetricsStats) error {
	h.Logger.Info("stats", zap.Any("stats", stats))
	return nil
}
