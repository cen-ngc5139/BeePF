package export

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// InternalBufferValueEventProcessor 内部缓冲区事件处理器
type InternalBufferValueEventProcessor interface {
	HandleEvent(data []byte) error
}

// InternalSampleMapProcessor 内部采样 Map 事件处理器
type InternalSampleMapProcessor interface {
	HandleEvent(data []byte) error
}

// JsonExportEventHandler JSON格式导出处理器
type JsonExportEventHandler struct {
	exporter *EventExporter
	mu       *sync.RWMutex
}

func NewJsonExportEventHandler(exporter *EventExporter) *JsonExportEventHandler {
	return &JsonExportEventHandler{
		exporter: exporter,
		mu:       &sync.RWMutex{},
	}
}

func (h *JsonExportEventHandler) HandleEvent(data []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 获取检查过的类型信息
	checkedTypes, err := h.exporter.InternalImpl.GetCheckedTypes()
	if err != nil {
		return fmt.Errorf("get checked types error: %w", err)
	}

	// 导出为JSON
	jsonData, err := DumpToJsonWithCheckedTypes(checkedTypes, data)
	if err != nil {
		return fmt.Errorf("dump to json error: %w", err)
	}

	// 输出数据
	return h.exporter.UserExportEventHandler.HandleEvent(h.exporter.UserCtx, &ReceivedEventData{
		Type:     TypeJsonText,
		JsonText: string(jsonData),
	})
}

// PlainTextExportEventHandler 纯文本导出处理器
type PlainTextExportEventHandler struct {
	exporter *EventExporter
	mu       *sync.RWMutex
}

func NewPlainTextExportEventHandler(exporter *EventExporter) *PlainTextExportEventHandler {
	return &PlainTextExportEventHandler{
		exporter: exporter,
		mu:       &sync.RWMutex{},
	}
}

func (h *PlainTextExportEventHandler) HandleEvent(data []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 获取检查过的类型信息
	checkedTypes, err := h.exporter.InternalImpl.GetCheckedTypes()
	if err != nil {
		return fmt.Errorf("get checked types error: %w", err)
	}

	// 生成输出
	var output strings.Builder

	// 添加时间戳
	now := time.Now().Format("15:04:05")
	fmt.Fprintf(&output, "%-8s ", now)

	// 导出数据
	err = DumpToStringWithCheckedTypes(checkedTypes, data, &output)
	if err != nil {
		return fmt.Errorf("dump to string error: %w", err)
	}

	// 输出数据
	return h.exporter.UserExportEventHandler.HandleEvent(h.exporter.UserCtx, &ReceivedEventData{
		Type: TypePlainText,
		Text: output.String(),
	})
}

// RawExportEventHandler 原始数据导出处理器
type RawExportEventHandler struct {
	exporter *EventExporter
	mu       *sync.RWMutex
}

func NewRawExportEventHandler(exporter *EventExporter) *RawExportEventHandler {
	return &RawExportEventHandler{
		exporter: exporter,
		mu:       &sync.RWMutex{},
	}
}

func (h *RawExportEventHandler) HandleEvent(data []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.exporter.UserExportEventHandler == nil {
		fmt.Println("Raw export event handler expects user callback, data will be dropped")
		return nil
	}

	// 直接传递原始数据
	return h.exporter.UserExportEventHandler.HandleEvent(h.exporter.UserCtx, &ReceivedEventData{
		Type:   TypeBuffer,
		Buffer: data,
	})
}
