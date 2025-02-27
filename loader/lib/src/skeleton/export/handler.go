package export

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/helper"
)

// InternalBufferValueEventProcessor 内部缓冲区事件处理器
type InternalBufferValueEventProcessor interface {
	HandleEvent(data []byte) error
}

// InternalSampleMapProcessor 内部采样 Map 事件处理器
type InternalSampleMapProcessor interface {
	HandleEvent(keyBuffer, valueBuffer []byte) error
}

// JsonExportEventHandler JSON格式导出处理器
type JsonExportEventHandler struct {
	Exporter *EventExporter
	Mu       *sync.RWMutex
}

func NewJsonExportEventHandler(exporter *EventExporter) *JsonExportEventHandler {
	return &JsonExportEventHandler{
		Exporter: exporter,
		Mu:       &sync.RWMutex{},
	}
}

func (h *JsonExportEventHandler) HandleEvent(data []byte) error {
	h.Mu.RLock()
	defer h.Mu.RUnlock()

	// 获取检查过的类型信息
	checkedTypes, err := h.Exporter.InternalImpl.GetCheckedTypes()
	if err != nil {
		return fmt.Errorf("get checked types error: %w", err)
	}

	// 导出为JSON
	jsonData, err := DumpToJsonWithCheckedTypes(checkedTypes, data)
	if err != nil {
		return fmt.Errorf("dump to json error: %w", err)
	}

	// 输出数据
	return h.Exporter.UserExportEventHandler.HandleEvent(h.Exporter.UserCtx, &ReceivedEventData{
		Type:     TypeJsonText,
		JsonText: string(jsonData),
	})
}

// PlainTextExportEventHandler 纯文本导出处理器
type PlainTextExportEventHandler struct {
	Exporter *EventExporter
	Mu       *sync.RWMutex
}

func NewPlainTextExportEventHandler(exporter *EventExporter) *PlainTextExportEventHandler {
	return &PlainTextExportEventHandler{
		Exporter: exporter,
		Mu:       &sync.RWMutex{},
	}
}

func (h *PlainTextExportEventHandler) HandleEvent(data []byte) error {
	h.Mu.RLock()
	defer h.Mu.RUnlock()

	// 获取检查过的类型信息
	checkedTypes, err := h.Exporter.InternalImpl.GetCheckedTypes()
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
	return h.Exporter.UserExportEventHandler.HandleEvent(h.Exporter.UserCtx, &ReceivedEventData{
		Type: TypePlainText,
		Text: output.String(),
	})
}

// RawExportEventHandler 原始数据导出处理器
type RawExportEventHandler struct {
	Exporter *EventExporter
	Mu       *sync.RWMutex
}

func NewRawExportEventHandler(exporter *EventExporter) *RawExportEventHandler {
	return &RawExportEventHandler{
		Exporter: exporter,
		Mu:       &sync.RWMutex{},
	}
}

func (h *RawExportEventHandler) HandleEvent(data []byte) error {
	h.Mu.RLock()
	defer h.Mu.RUnlock()

	if h.Exporter.UserExportEventHandler == nil {
		fmt.Println("Raw export event handler expects user callback, data will be dropped")
		return nil
	}

	// 直接传递原始数据
	return h.Exporter.UserExportEventHandler.HandleEvent(h.Exporter.UserCtx, &ReceivedEventData{
		Type:   TypeBuffer,
		Buffer: data,
	})
}

// JsonMapExporter JSON 格式导出处理器
type JsonMapExporter struct {
	Exporter *EventExporter
}

func NewJsonMapExporter(exporter *EventExporter) *JsonMapExporter {
	return &JsonMapExporter{Exporter: exporter}
}

func (h *JsonMapExporter) HandleEvent(keyBuffer, valueBuffer []byte) error {
	// 获取检查过的类型信息
	checkedKeyTypes, err := h.Exporter.InternalImpl.GetCheckedKeyTypes()
	if err != nil {
		return fmt.Errorf("get checked types error: %w", err)
	}

	checkedValueTypes, err := h.Exporter.InternalImpl.GetCheckedValueTypes()
	if err != nil {
		return fmt.Errorf("get checked types error: %w", err)
	}

	// 导出 key
	keyOut, err := DumpToJsonWithCheckedTypes(checkedKeyTypes, keyBuffer)
	if err != nil {
		return fmt.Errorf("dump key to json error: %w", err)
	}

	// 导出 value
	valueOut, err := DumpToJsonWithCheckedTypes(checkedValueTypes, valueBuffer)
	if err != nil {
		return fmt.Errorf("dump value to json error: %w", err)
	}

	// 构造最终的 JSON
	result := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"key":       keyOut,
		"value":     valueOut,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal json error: %w", err)
	}

	// 输出数据
	h.Exporter.UserExportEventHandler.HandleEvent(h.Exporter.UserCtx, &ReceivedEventData{
		Type:     TypeJsonText,
		JsonText: string(jsonData),
	})

	return nil
}

// PlainTextMapExporter 纯文本导出处理器
type PlainTextMapExporter struct {
	Exporter *EventExporter
}

func NewPlainTextMapExporter(exporter *EventExporter) *PlainTextMapExporter {
	return &PlainTextMapExporter{Exporter: exporter}
}

func (h *PlainTextMapExporter) HandleEvent(keyBuffer, valueBuffer []byte) error {
	// 获取检查过的类型信息
	checkedKeyTypes, err := h.Exporter.InternalImpl.GetCheckedKeyTypes()
	if err != nil {
		return fmt.Errorf("get checked types error: %w", err)
	}

	checkedValueTypes, err := h.Exporter.InternalImpl.GetCheckedValueTypes()
	if err != nil {
		return fmt.Errorf("get checked types error: %w", err)
	}

	var output strings.Builder
	output.WriteString("key = ")

	// 导出 key
	if err := DumpToStringWithCheckedTypes(checkedKeyTypes, keyBuffer, &output); err != nil {
		return fmt.Errorf("dump key error: %w", err)
	}
	output.WriteString("\n")

	// 导出 value
	if err := DumpToStringWithCheckedTypes(checkedValueTypes, valueBuffer, &output); err != nil {
		return fmt.Errorf("dump value error: %w", err)
	}

	// 输出数据
	h.Exporter.UserExportEventHandler.HandleEvent(h.Exporter.UserCtx, &ReceivedEventData{
		Type: TypePlainText,
		Text: output.String(),
	})

	return nil
}

// RawMapExporter 原始数据导出处理器
type RawMapExporter struct {
	Exporter *EventExporter
}

func NewRawMapExporter(exporter *EventExporter) *RawMapExporter {
	return &RawMapExporter{Exporter: exporter}
}

func (h *RawMapExporter) HandleEvent(keyBuffer, valueBuffer []byte) error {
	// 直接传递原始数据
	return h.Exporter.UserExportEventHandler.HandleEvent(h.Exporter.UserCtx, &ReceivedEventData{
		Type:   TypeBuffer,
		Buffer: append(keyBuffer, valueBuffer...),
	})
}

// Log2HistExporter 直方图导出处理器
type Log2HistExporter struct {
	Exporter *EventExporter
}

func NewLog2HistExporter(exporter *EventExporter) *Log2HistExporter {
	return &Log2HistExporter{Exporter: exporter}
}

func (h *Log2HistExporter) HandleEvent(keyBuffer, valueBuffer []byte) error {
	// 获取检查过的类型信息
	checkedKeyTypes, err := h.Exporter.InternalImpl.GetCheckedKeyTypes()
	if err != nil {
		return fmt.Errorf("get map config error: %w", err)
	}

	checkedValueTypes, err := h.Exporter.InternalImpl.GetCheckedValueTypes()
	if err != nil {
		return fmt.Errorf("get map config error: %w", err)
	}

	var output strings.Builder
	output.WriteString("key = ")

	// 导出 key
	if err := DumpToStringWithCheckedTypes(checkedKeyTypes, keyBuffer, &output); err != nil {
		return fmt.Errorf("dump key error: %w", err)
	}
	output.WriteString("\n")

	// 查找并处理 slots
	var slots []uint32
	for _, member := range checkedValueTypes {
		offset := member.BitOffset / 8
		if member.BitOffset%8 != 0 {
			return fmt.Errorf("bit fields not supported")
		}

		if member.FieldName == "slots" {
			// 读取 slots 数据
			slotCount := member.Size / 4
			slots = make([]uint32, slotCount)
			for i := uint32(0); i < uint32(slotCount); i++ {
				start := offset.Bytes() + i*4
				end := start + 4
				slots[i] = binary.LittleEndian.Uint32(valueBuffer[start:end])
			}
		} else {
			// 导出其他字段
			output.WriteString(fmt.Sprintf("%s = ", member.FieldName))
			val, err := DumpToString(member.Type, valueBuffer[offset:offset+member.Size])
			if err != nil {
				return fmt.Errorf("dump value field error: %w", err)
			}
			output.WriteString(val)
			output.WriteString("\n")
		}
	}

	// 输出基本信息
	h.Exporter.UserExportEventHandler.HandleEvent(h.Exporter.UserCtx, &ReceivedEventData{
		Type: TypePlainText,
		Text: output.String(),
	})

	// 如果有 slots，打印直方图
	if len(slots) > 0 {

		h.Exporter.UserExportEventHandler.HandleEvent(h.Exporter.UserCtx, &ReceivedEventData{
			Type: TypePlainText,
			// TODO: 从命令传入直方图的字段名称
			Text: helper.PrintLog2Hist(slots, ""),
		})
	}

	return nil
}
