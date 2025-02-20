package loader

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/container"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/ringbuf"
	"go.uber.org/zap"
)

// Loader 定义加载器接口
type Loader interface {
	Init() error
	Load() error
	Start() error
	Stop() error
}

// MapHandler 定义 Map 处理器接口
type MapHandler interface {
	Type() ebpf.MapType
	Setup(*ebpf.Map) (*skeleton.ProgramPoller, error)
}

// BaseMapHandler 提供通用实现
type BaseMapHandler struct {
	logger       *zap.Logger
	config       *Config
	collection   *ebpf.Collection
	btfContainer *container.BTFContainer
}

// setupExporter 设置事件导出器
func (h *BaseMapHandler) setupExporter(structType *btf.Struct) (*export.EventExporter, error) {
	ee := export.NewEventExporterBuilder().
		SetExportFormat(export.FormatJson).
		SetUserContext(export.NewUserContext(0)).
		SetEventHandler(&export.MyCustomHandler{Logger: h.logger})

	exporter, err := ee.BuildForSingleValueWithTypeDescriptor(
		&export.BTFTypeDescriptor{
			Type: structType,
			Name: structType.TypeName(),
		},
		h.btfContainer,
	)
	if err != nil {
		return nil, fmt.Errorf("build event exporter failed: %w", err)
	}

	return exporter, nil
}

// setupPoller 设置轮询器
func (h *BaseMapHandler) setupPoller(poller skeleton.Poller) (*skeleton.ProgramPoller, error) {
	// 创建程序轮询器
	programPoller := skeleton.NewProgramPoller(h.config.PollTimeout)

	// 启动轮询
	go func() {
		programPoller.StartPolling(
			h.config.ProgramName,
			poller.GetPollFunc(),
			h.handlePollingError,
		)
	}()

	return programPoller, nil
}

// findTargetStruct 查找目标结构体
func (h *BaseMapHandler) findTargetStruct() (*btf.Struct, error) {
	for _, v := range h.collection.Variables {
		structType, err := skeleton.FindStructType(v.Type())
		if err != nil {
			h.logger.Warn("find struct type failed", zap.Error(err))
			continue
		}

		if structType.Name == h.config.StructName {
			return structType, nil
		}
	}
	return nil, fmt.Errorf("target struct %s not found", h.config.StructName)
}

// PerfEventMapHandler 实现
type PerfEventMapHandler struct {
	BaseMapHandler
}

func (h *PerfEventMapHandler) Type() ebpf.MapType {
	return ebpf.PerfEventArray
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
		Timeout:   h.config.PollTimeout,
	}

	// 设置轮询器
	return h.setupPoller(poller)
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
		Timeout:   h.config.PollTimeout,
	}

	return h.setupPoller(poller)
}

// BPFLoader 实现 eBPF 程序加载器
type BPFLoader struct {
	logger          *zap.Logger
	config          *Config
	collection      *ebpf.Collection
	skeleton        *skeleton.BpfSkeleton
	preLoadSkeleton *skeleton.PreLoadBpfSkeleton
	pollers         []*skeleton.ProgramPoller
	eventExporter   *export.EventExporter
	btfContainer    *container.BTFContainer
	mapHandlers     []MapHandler
}

// Config 配置结构
type Config struct {
	ObjectPath   string
	BTFPath      string
	Logger       *zap.Logger
	PollTimeout  time.Duration
	ProgramName  string
	StructName   string
	BTFContainer *container.BTFContainer
}

func NewBPFLoader(cfg *Config) *BPFLoader {
	loader := &BPFLoader{
		logger:      cfg.Logger,
		config:      cfg,
		mapHandlers: make([]MapHandler, 0),
	}

	// 注册默认的 map 处理器
	loader.RegisterMapHandler(&PerfEventMapHandler{
		logger:       cfg.Logger,
		config:       cfg,
		btfContainer: cfg.BTFContainer,
	})

	loader.RegisterMapHandler(&RingBufMapHandler{
		logger:       cfg.Logger,
		config:       cfg,
		btfContainer: cfg.BTFContainer,
	})

	return loader
}

// Init 初始化阶段
func (l *BPFLoader) Init() error {
	l.logger.Info("initializing BPF loader...")

	// 生成组合对象
	pkg, err := meta.GenerateComposedObject(l.config.ObjectPath)
	if err != nil {
		return fmt.Errorf("generate composed object failed: %w", err)
	}

	// 构建预加载骨架
	preLoadSkeleton, err := skeleton.FromJsonPackage(pkg, filepath.Dir(l.config.ObjectPath)).Build()
	if err != nil {
		return fmt.Errorf("build preload skeleton failed: %w", err)
	}

	l.preLoadSkeleton = preLoadSkeleton
	return nil
}

// Load 加载阶段
func (l *BPFLoader) Load() error {
	l.logger.Info("loading BPF programs...")

	// 加载并附加 eBPF 程序
	skel, err := l.preLoadSkeleton.LoadAndAttach()
	if err != nil {
		return fmt.Errorf("load and attach BPF programs failed: %w", err)
	}

	l.collection = skel.Collection
	return nil
}

// RegisterMapHandler 注册 Map 处理器
func (l *BPFLoader) RegisterMapHandler(handler MapHandler) {
	l.mapHandlers = append(l.mapHandlers, handler)
}

// Start 启动阶段
func (l *BPFLoader) Start() error {
	l.logger.Info("starting BPF programs...")

	// 处理所有 maps
	for _, m := range l.collection.Maps {
		// 查找对应的处理器
		for _, handler := range l.mapHandlers {
			if handler.Type() == m.Type() {
				poller, err := handler.Setup(m)
				if err != nil {
					return fmt.Errorf("setup map handler failed: %w", err)
				}
				if poller != nil {
					l.pollers = append(l.pollers, poller)
				}
				break
			}
		}
	}

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		l.logger.Info("received signal, stopping...")
		l.Stop()
	}()

	return nil
}

// Stop 停止阶段
func (l *BPFLoader) Stop() error {
	l.logger.Info("stopping BPF programs...")

	// 停止所有 poller
	for _, p := range l.pollers {
		p.Stop()
	}

	// 关闭集合
	if l.collection != nil {
		l.collection.Close()
	}

	return nil
}

// handlePollingError 处理轮询错误
func (h *BaseMapHandler) handlePollingError(err error) {
	h.logger.Error("polling error", zap.Error(err))
}
