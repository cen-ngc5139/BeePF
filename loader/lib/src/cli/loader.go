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
	"github.com/cilium/ebpf/link"
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

// BPFLoader 实现 eBPF 程序加载器
type BPFLoader struct {
	Logger          *zap.Logger
	Config          *Config
	Collection      *ebpf.Collection
	Skeleton        *skeleton.BpfSkeleton
	PreLoadSkeleton *skeleton.PreLoadBpfSkeleton
	Pollers         []*skeleton.ProgramPoller
	MapHandlers     []MapHandler
	BTFContainer    *container.BTFContainer
	Links           []link.Link
	done            chan struct{}
}

// Config 配置结构
type Config struct {
	ObjectPath  string
	BTFPath     string
	Logger      *zap.Logger
	PollTimeout time.Duration
	ProgramName string
	StructName  string
}

func NewBPFLoader(cfg *Config) *BPFLoader {
	loader := &BPFLoader{
		Logger:      cfg.Logger,
		Config:      cfg,
		MapHandlers: make([]MapHandler, 0),
	}

	// 注册默认的 map 处理器
	loader.RegisterMapHandler(&PerfEventMapHandler{
		BaseMapHandler: BaseMapHandler{
			Logger: cfg.Logger,
			Config: cfg,
		},
	})

	loader.RegisterMapHandler(&RingBufMapHandler{
		BaseMapHandler: BaseMapHandler{
			Logger: cfg.Logger,
			Config: cfg,
		},
	})

	return loader
}

// Init 初始化阶段
func (l *BPFLoader) Init() error {
	l.Logger.Info("initializing BPF loader...")

	// 生成组合对象
	pkg, err := meta.GenerateComposedObject(l.Config.ObjectPath)
	if err != nil {
		return fmt.Errorf("generate composed object failed: %w", err)
	}

	// 构建预加载骨架
	preLoadSkeleton, err := skeleton.FromJsonPackage(pkg, filepath.Dir(l.Config.ObjectPath)).Build()
	if err != nil {
		return fmt.Errorf("build preload skeleton failed: %w", err)
	}

	l.PreLoadSkeleton = preLoadSkeleton
	return nil
}

// Load 加载阶段
func (l *BPFLoader) Load() error {
	l.Logger.Info("loading BPF programs...")

	// 加载并附加 eBPF 程序
	skel, err := l.PreLoadSkeleton.LoadAndAttach()
	if err != nil {
		return fmt.Errorf("load and attach BPF programs failed: %w", err)
	}

	l.Collection = skel.Collection
	l.BTFContainer = skel.Btf
	l.Links = skel.Links
	for _, handler := range l.MapHandlers {
		handler.SetCollection(l.Collection)
		handler.SetBTFContainer(l.BTFContainer)
	}
	return nil
}

// RegisterMapHandler 注册 Map 处理器
func (l *BPFLoader) RegisterMapHandler(handler MapHandler) {
	l.MapHandlers = append(l.MapHandlers, handler)
}

// Start 启动阶段
func (l *BPFLoader) Start() error {
	l.Logger.Info("starting BPF programs...")

	// 处理所有 maps
	for _, m := range l.Collection.Maps {
		// 查找对应的处理器
		for _, handler := range l.MapHandlers {
			if handler.Type() == m.Type() {
				poller, err := handler.Setup(m)
				if err != nil {
					return fmt.Errorf("setup map handler failed: %w", err)
				}
				if poller != nil {
					l.Pollers = append(l.Pollers, poller)
				}
				break
			}
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建 done channel
	l.done = make(chan struct{})

	go func() {
		sig := <-sigChan
		l.Logger.Info("received signal", zap.String("signal", sig.String()))

		if err := l.Stop(); err != nil {
			l.Logger.Error("stop failed", zap.Error(err))
		}

		// 通知主程序可以退出
		close(l.done)

		// 停止接收更多信号
		signal.Stop(sigChan)
	}()

	return nil
}

// Stop 停止阶段
func (l *BPFLoader) Stop() error {
	l.Logger.Info("starting cleanup process...")

	// 1. 先停止所有 poller，因为它们在使用 maps
	l.Logger.Info("stopping pollers")
	for _, p := range l.Pollers {
		p.Stop()
	}

	// 2. 关闭所有 map handlers，它们持有 map readers
	l.Logger.Info("closing map handlers")
	for _, handler := range l.MapHandlers {
		handler.Close()
	}

	// 3. 关闭所有 links，因为它们引用了 programs
	l.Logger.Info("closing links")
	for _, link := range l.Links {
		if err := link.Close(); err != nil {
			l.Logger.Error("failed to close link", zap.Error(err))
		}
	}

	// 4. 最后关闭 collection，它会关闭所有 maps 和 programs
	l.Logger.Info("closing collection")
	if l.Collection != nil {
		l.Collection.Close()
	}

	if l.Config.Logger != nil {
		l.Config.Logger.Info("closing logger")
		l.Config.Logger.Sync()
	}

	// 5. 清空所有引用
	l.Pollers = nil
	l.MapHandlers = nil
	l.Links = nil
	l.Collection = nil

	return nil
}

// handlePollingError 处理轮询错误
func (h *BaseMapHandler) handlePollingError(err error) {
	h.Logger.Error("polling error", zap.Error(err))
}

// Done 返回完成信号通道
func (l *BPFLoader) Done() <-chan struct{} {
	return l.done
}
