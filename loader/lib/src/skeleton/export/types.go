package export

import (
	"github.com/cen-ngc5139/BeePF/loader/lib/src/container"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cilium/ebpf/btf"
)

// ExportFormatType 定义导出格式类型
type ExportFormatType int

const (
	FormatJson ExportFormatType = iota
	FormatPlainText
	FormatRawEvent
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
	HandleEvent(ctx interface{}, data *ReceivedEventData) error
}

// EventExporter 主要的导出器结构
type EventExporter struct {
	userExportEventHandler EventHandler
	userCtx                interface{}
	btfContainer           *container.BTFContainer
	internalImpl           ExporterImplementation
}

// ExporterImplementation 定义内部实现
type ExporterImplementation struct {
	Type            ImplType
	EventProcessor  interface{} // 可以是 BufferValueProcessor 或 SampleMapProcessor
	CheckedTypes    []CheckedExportedMember
	CheckedKeyTypes []CheckedExportedMember
	CheckedValTypes []CheckedExportedMember
	SampleMapConfig *meta.MapSampleMeta
}

type ImplType int

const (
	ImplBufferValue ImplType = iota
	ImplKeyValueMap
)

// EventExporterBuilder 用于构建 EventExporter
type EventExporterBuilder struct {
	exportFormat       ExportFormatType
	exportEventHandler EventHandler
	userCtx            interface{}
}
