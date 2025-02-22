package loader

import (
	"fmt"
	"os"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/container"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/ringbuf"
	"go.uber.org/zap"
)

// MapHandler 定义 Map 处理器接口
type MapHandler interface {
	Type() ebpf.MapType
	Setup(*ebpf.Map) (*skeleton.ProgramPoller, error)
	SetCollection(*ebpf.Collection)
	SetBTFContainer(*container.BTFContainer)
	Close()
}

// BaseMapHandler 提供通用实现
type BaseMapHandler struct {
	Logger       *zap.Logger
	Config       *Config
	Collection   *ebpf.Collection
	BTFContainer *container.BTFContainer
	Poller       skeleton.Poller
	Stats        *metrics.Collector
}

// setupExporter 设置事件导出器
func (h *BaseMapHandler) setupExporter(structType *btf.Struct) (*export.EventExporter, error) {
	ee := export.NewEventExporterBuilder().
		SetExportFormat(export.FormatJson).
		SetUserContext(export.NewUserContext(0)).
		SetEventHandler(&export.MyCustomHandler{Logger: h.Logger})

	exporter, err := ee.BuildForSingleValueWithTypeDescriptor(
		&export.BTFTypeDescriptor{
			Type: structType,
			Name: structType.TypeName(),
		},
		h.BTFContainer,
	)
	if err != nil {
		return nil, fmt.Errorf("build event exporter failed: %w", err)
	}

	return exporter, nil
}

// setupPoller 设置轮询器
func (h *BaseMapHandler) setupPoller(poller skeleton.Poller) (*skeleton.ProgramPoller, error) {
	h.Poller = poller
	// 创建程序轮询器
	programPoller := skeleton.NewProgramPoller(h.Config.PollTimeout)

	// 启动轮询

	programPoller.StartPolling(
		h.Config.ProgramName,
		poller.GetPollFunc(),
		h.handlePollingError,
	)

	return programPoller, nil
}

// findTargetStruct 查找目标结构体
func (h *BaseMapHandler) findTargetStruct() (*btf.Struct, error) {
	for _, v := range h.Collection.Variables {
		structType, err := skeleton.FindStructType(v.Type())
		if err != nil {
			h.Logger.Warn("find struct type failed", zap.Error(err))
			continue
		}

		if structType.Name == h.Config.StructName {
			return structType, nil
		}
	}
	return nil, fmt.Errorf("target struct %s not found", h.Config.StructName)
}

// PerfEventMapHandler 实现
type PerfEventMapHandler struct {
	BaseMapHandler
}

func (h *PerfEventMapHandler) Type() ebpf.MapType {
	return ebpf.PerfEventArray
}

func (h *PerfEventMapHandler) SetCollection(collection *ebpf.Collection) {
	h.Collection = collection
}

func (h *PerfEventMapHandler) SetBTFContainer(btfContainer *container.BTFContainer) {
	h.BTFContainer = btfContainer
}

func (h *PerfEventMapHandler) Setup(m *ebpf.Map) (*skeleton.ProgramPoller, error) {
	// 创建读取器
	reader, err := perf.NewReader(m, os.Getpagesize())
	if err != nil {
		return nil, fmt.Errorf("create perf reader failed: %w", err)
	}

	// 查找目标结构体
	structType, err := h.findTargetStruct()
	if err != nil {
		return nil, err
	}

	// 设置导出器
	exporter, err := h.setupExporter(structType)
	if err != nil {
		return nil, err
	}

	// 创建处理器
	processor := export.NewJsonExportEventHandler(exporter)

	poller := &skeleton.PerfEventPoller{
		Reader:    reader,
		Processor: processor,
		Timeout:   h.Config.PollTimeout,
	}

	// 设置轮询器
	return h.setupPoller(poller)
}

func (h *PerfEventMapHandler) Close() {
	if h.Poller != nil {
		h.Poller.Close()
	}
}

// RingBufMapHandler 实现
type RingBufMapHandler struct {
	BaseMapHandler
}

func (h *RingBufMapHandler) Type() ebpf.MapType {
	return ebpf.RingBuf
}

func (h *RingBufMapHandler) Setup(m *ebpf.Map) (*skeleton.ProgramPoller, error) {
	// 创建读取器
	reader, err := ringbuf.NewReader(m)
	if err != nil {
		return nil, fmt.Errorf("create ring buffer reader failed: %w", err)
	}

	// 使用相同的通用逻辑
	structType, err := h.findTargetStruct()
	if err != nil {
		return nil, err
	}

	exporter, err := h.setupExporter(structType)
	if err != nil {
		return nil, err
	}

	processor := export.NewJsonExportEventHandler(exporter)
	poller := &skeleton.RingBufPoller{
		Reader:    reader,
		Processor: processor,
		Timeout:   h.Config.PollTimeout,
	}

	return h.setupPoller(poller)
}

func (h *RingBufMapHandler) Close() {
	if h.Poller != nil {
		h.Poller.Close()
	}
}

func (h *RingBufMapHandler) SetCollection(collection *ebpf.Collection) {
	h.Collection = collection
}

func (h *RingBufMapHandler) SetBTFContainer(btfContainer *container.BTFContainer) {
	h.BTFContainer = btfContainer
}
