package metrics

import "go.uber.org/zap"

type Handler interface {
	Handle(stats *Stats) error
}

type DefaultHandler struct {
	Logger *zap.Logger
}

func NewDefaultHandler(logger *zap.Logger) *DefaultHandler {
	return &DefaultHandler{
		Logger: logger,
	}
}

func (h *DefaultHandler) Handle(stats *Stats) error {
	h.Logger.Info("stats", zap.Any("stats", stats))
	return nil
}
