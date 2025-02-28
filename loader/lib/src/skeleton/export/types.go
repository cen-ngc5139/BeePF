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
	FormatLog2Hist
)

// 为了兼容性，保留类型别名
type DataType = meta.DataType
type ReceivedEventData = meta.ReceivedEventData
type EventHandler = meta.EventHandler
type UserContext = meta.UserContext

// 常量别名
const (
	TypeBuffer         = meta.TypeBuffer
	TypeKeyValueBuffer = meta.TypeKeyValueBuffer
	TypePlainText      = meta.TypePlainText
	TypeJsonText       = meta.TypeJsonText
)

// CheckedExportedMember 表示已检查的导出结构成员
type CheckedExportedMember struct {
	FieldName          string
	Type               btf.Type
	BitOffset          btf.Bits
	Size               btf.Bits
	OutputHeaderOffset btf.Bits
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

// 使用 meta 包中的函数
var NewUserContext = meta.NewUserContext
