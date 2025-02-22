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
	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"go.uber.org/zap"
)

// Loader 定义加载器接口
type Loader interface {
	Init() error
	Load() error
	Start() error
	Stop() error
	Stats() error
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
	StatsCollector  metrics.Collector
}

// Config 配置结构
type Config struct {
	ObjectPath    string
	BTFPath       string
	Logger        *zap.Logger
	PollTimeout   time.Duration
	ProgramName   string
	StructName    string
	IsEnableStats bool
	StatsInterval time.Duration
}

func NewBPFLoader(cfg *Config) *BPFLoader {
	loader := &BPFLoader{
		Logger:      cfg.Logger,
		Config:      cfg,
		MapHandlers: make([]MapHandler, 0),
	}

	if cfg.IsEnableStats {
		if cfg.StatsInterval == 0 {
			cfg.StatsInterval = 1 * time.Second
		}

		collector, err := metrics.NewStatsCollector(cfg.StatsInterval)
		if err != nil {
			cfg.Logger.Error("failed to create stats collector", zap.Error(err))
		}
		loader.StatsCollector = collector
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

	l.Logger.Info("stopping stats collector")
	if l.StatsCollector != nil {
		err := l.StatsCollector.Stop()
		if err != nil {
			l.Logger.Error("failed to stop stats collector", zap.Error(err))
		}
	}

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

func (l *BPFLoader) Stats() error {
	l.Logger.Info("collecting stats")

	if l.StatsCollector == nil {
		l.Logger.Info("stats collector is not enabled")
		return nil
	}

	attachedPros := make(map[uint32]*ebpf.Program)
	for _, prog := range l.Collection.Programs {
		info, err := prog.Info()
		if err != nil {
			return fmt.Errorf("failed to get program info: %w", err)
		}

		id, ok := info.ID()
		if !ok {
			return fmt.Errorf("failed to get program id: %w", err)
		}
		attachedPros[uint32(id)] = prog
	}

	if err := l.StatsCollector.SetAttachedPros(attachedPros); err != nil {
		return err
	}

	return l.StatsCollector.Start()
}

// handlePollingError 处理轮询错误
func (h *BaseMapHandler) handlePollingError(err error) {
	h.Logger.Error("polling error", zap.Error(err))
}

// Done 返回完成信号通道
func (l *BPFLoader) Done() <-chan struct{} {
	return l.done
}
