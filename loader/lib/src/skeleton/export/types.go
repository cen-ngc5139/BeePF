package export

import (
	"sync"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/container"
	"github.com/cilium/ebpf/btf"
)

// ExportFormatType 定义导出格式类型
type ExportFormatType int

const (
	FormatJson ExportFormatType = iota
	FormatPlainText
	FormatRawEvent
	FormatLog2Hist
)

// ReceivedEventData 表示从 eBPF 程序接收到的数据
type ReceivedEventData struct {
	Type     DataType
	Buffer   []byte
	KeyBuf   []byte
	ValueBuf []byte
	Text     string
	JsonText string
}

type DataType int

const (
	TypeBuffer DataType = iota
	TypeKeyValueBuffer
	TypePlainText
	TypeJsonText
)

// CheckedExportedMember 表示已检查的导出结构成员
type CheckedExportedMember struct {
	FieldName          string
	Type               btf.Type
	BitOffset          btf.Bits
	Size               btf.Bits
	OutputHeaderOffset btf.Bits
}

// EventHandler 定义事件处理接口
type EventHandler interface {
	HandleEvent(ctx *UserContext, data *ReceivedEventData) error
}

// UserContext 用于存储用户上下文
type UserContext struct {
	Value interface{}
	mu    sync.RWMutex
}

// EventExporter 主要的导出器结构
type EventExporter struct {
	UserExportEventHandler EventHandler
	UserCtx                *UserContext
	BTFContainer           *container.BTFContainer
	InternalImpl           ExporterInternalImplementation
}

// EventExporterBuilder 用于构建 EventExporter
type EventExporterBuilder struct {
	ExportFormat       ExportFormatType
	ExportEventHandler EventHandler
	UserCtx            *UserContext
}

// NewUserContext 创建新的用户上下文
func NewUserContext(value interface{}) *UserContext {
	return &UserContext{
		Value: value,
	}
}

// GetValue 获取上下文值
func (c *UserContext) GetValue() interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Value
}

// SetValue 设置上下文值
func (c *UserContext) SetValue(value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Value = value
}
