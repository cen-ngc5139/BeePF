package export

import (
	"fmt"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/container"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
)

// NewEventExporterBuilder 创建新的构建器
func NewEventExporterBuilder() *EventExporterBuilder {
	return &EventExporterBuilder{
		ExportFormat: FormatPlainText, // 默认纯文本格式
	}
}

// SetExportFormat 设置导出格式
func (b *EventExporterBuilder) SetExportFormat(format ExportFormatType) *EventExporterBuilder {
	b.ExportFormat = format
	return b
}

// SetEventHandler 设置事件处理器
func (b *EventExporterBuilder) SetEventHandler(handler EventHandler) *EventExporterBuilder {
	b.ExportEventHandler = handler
	return b
}

// SetUserContext 设置用户上下文
func (b *EventExporterBuilder) SetUserContext(ctx *UserContext) *EventExporterBuilder {
	b.UserCtx = ctx
	return b
}

func (b *EventExporterBuilder) BuildForSingleValueWithTypeDescriptor(
	typeDesc TypeDescriptor,
	btfContainer *container.BTFContainer,
) (*EventExporter, error) {
	// 1. 参数验证
	if btfContainer == nil {
		return nil, fmt.Errorf("BTF container is required")
	}
	if typeDesc == nil {
		return nil, fmt.Errorf("type descriptor is required")
	}

	// 创建 EventExporter
	exporter := &EventExporter{
		BTFContainer:           btfContainer,
		UserExportEventHandler: b.ExportEventHandler,
		UserCtx:                b.UserCtx,
	}

	checkedTypes, err := typeDesc.BuildCheckedExportedMembers()
	if err != nil {
		return nil, fmt.Errorf("failed to build checked exported members: %w", err)
	}

	// 3. 创建内部处理器
	var processor InternalBufferValueEventProcessor
	switch b.ExportFormat {
	case FormatJson:
		processor = NewJsonExportEventHandler(exporter)
	case FormatPlainText:
		processor = NewPlainTextExportEventHandler(exporter)
	case FormatRawEvent:
		processor = NewRawExportEventHandler(exporter)
	default:
		return nil, fmt.Errorf("unsupported export format: %v", b.ExportFormat)
	}

	// 6. 创建导出器
	exporter.InternalImpl = &BufferValueProcessor{
		Processor:    processor,
		CheckedTypes: checkedTypes,
	}

	return exporter, nil
}

// BuildForSingleValue 构建用于单个值的导出器
func (b *EventExporterBuilder) BuildForSingleValue(
	exportType *meta.ExportedTypesStructMeta,
	btfContainer *container.BTFContainer,
	interpreter *meta.BufferValueInterpreter,
) (*EventExporter, error) {
	btfTypeDesc := &BTFTypeDescriptor{
		Type: exportType.Type,
		Name: exportType.Name,
	}

	return b.BuildForSingleValueWithTypeDescriptor(
		btfTypeDesc,
		btfContainer,
	)
}

// BuildForKeyValue 构建用于 key-value 的导出器
func (b *EventExporterBuilder) BuildForKeyValue(
	sampleConfig *meta.MapSampleMeta,
	exportKeyType *meta.ExportedTypesStructMeta,
	exportValueType *meta.ExportedTypesStructMeta,
	btfContainer *container.BTFContainer,
) (*EventExporter, error) {
	keyTypeDesc := &BTFTypeDescriptor{
		Type: exportKeyType.Type,
		Name: exportKeyType.Name,
	}
	valueTypeDesc := &BTFTypeDescriptor{
		Type: exportValueType.Type,
		Name: exportValueType.Name,
	}

	return b.BuildForKeyValueWithTypeDesc(
		keyTypeDesc,
		valueTypeDesc,
		btfContainer,
		sampleConfig,
	)
}

func (b *EventExporterBuilder) BuildForKeyValueWithTypeDesc(
	keyTypeDesc TypeDescriptor,
	valueTypeDesc TypeDescriptor,
	btfContainer *container.BTFContainer,
	sampleConfig *meta.MapSampleMeta,
) (*EventExporter, error) {
	// 1. 参数验证
	if btfContainer == nil {
		return nil, fmt.Errorf("BTF container is required")
	}
	if keyTypeDesc == nil {
		return nil, fmt.Errorf("key type descriptor is required")
	}
	if valueTypeDesc == nil {
		return nil, fmt.Errorf("value type descriptor is required")
	}

	// 2. 构建已检查的导出成员
	keyCheckedTypes, err := keyTypeDesc.BuildCheckedExportedMembers()
	if err != nil {
		return nil, fmt.Errorf("build key checked members error: %w", err)
	}

	valueCheckedTypes, err := valueTypeDesc.BuildCheckedExportedMembers()
	if err != nil {
		return nil, fmt.Errorf("build value checked members error: %w", err)
	}

	// 创建 EventExporter
	exporter := &EventExporter{
		BTFContainer:           btfContainer,
		UserExportEventHandler: b.ExportEventHandler,
		UserCtx:                b.UserCtx,
	}

	// 3. 创建内部处理器
	var processor InternalSampleMapProcessor
	switch b.ExportFormat {
	case FormatJson:
		processor = NewJsonMapExporter(exporter)
	case FormatPlainText:
		processor = NewPlainTextMapExporter(exporter)
	case FormatRawEvent:
		processor = NewRawMapExporter(exporter)
	case FormatLog2Hist:
		processor = NewLog2HistExporter(exporter)
	default:
		return nil, fmt.Errorf("unsupported export format: %v", b.ExportFormat)
	}

	// 4. 创建导出器
	exporter.InternalImpl = &KeyValueMapProcessor{
		Processor:         processor,
		CheckedKeyTypes:   keyCheckedTypes,
		CheckedValueTypes: valueCheckedTypes,
		MapConfig:         sampleConfig,
	}

	return exporter, nil
}
