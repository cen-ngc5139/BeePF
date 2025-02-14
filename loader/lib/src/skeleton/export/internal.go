// pkg/export_event/exporter.go

package export

import (
	"errors"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
)

// ExporterInternalImplementation 导出器内部实现类型
type ExporterInternalImplementation interface {
	GetCheckedTypes() ([]CheckedExportedMember, error)
	GetCheckedKeyTypes() ([]CheckedExportedMember, error)
	GetCheckedValueTypes() ([]CheckedExportedMember, error)
	Process(data []byte) error
}

// BufferValueProcessor Buffer 处理器实现
type BufferValueProcessor struct {
	Processor    InternalBufferValueEventProcessor
	CheckedTypes []CheckedExportedMember
}

func (b *BufferValueProcessor) GetCheckedTypes() ([]CheckedExportedMember, error) {
	if b.CheckedTypes == nil {
		return nil, errors.New("checked types is nil")
	}

	return b.CheckedTypes, nil
}

func (b *BufferValueProcessor) GetCheckedKeyTypes() ([]CheckedExportedMember, error) {
	return nil, errors.New("buffer value processor does not support get checked key types")
}

func (b *BufferValueProcessor) GetCheckedValueTypes() ([]CheckedExportedMember, error) {
	return nil, errors.New("buffer value processor does not support get checked value types")
}

func (b *BufferValueProcessor) Process(data []byte) error {
	return b.Processor.HandleEvent(data)
}

// KeyValueMapProcessor Map 处理器实现
type KeyValueMapProcessor struct {
	Processor         InternalSampleMapProcessor
	CheckedKeyTypes   []CheckedExportedMember
	CheckedValueTypes []CheckedExportedMember
	MapConfig         *meta.MapSampleMeta
}

func (k *KeyValueMapProcessor) GetCheckedTypes() ([]CheckedExportedMember, error) {
	return nil, errors.New("key value map processor does not support get checked types")
}

func (k *KeyValueMapProcessor) GetCheckedKeyTypes() ([]CheckedExportedMember, error) {
	if k.CheckedKeyTypes == nil {
		return nil, errors.New("checked key types is nil")
	}

	return k.CheckedKeyTypes, nil
}

func (k *KeyValueMapProcessor) GetCheckedValueTypes() ([]CheckedExportedMember, error) {
	if k.CheckedValueTypes == nil {
		return nil, errors.New("checked value types is nil")
	}
	return k.CheckedValueTypes, nil
}

func (k *KeyValueMapProcessor) Process(data []byte) error {
	return k.Processor.HandleEvent(data)
}
