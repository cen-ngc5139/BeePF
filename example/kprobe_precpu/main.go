package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	meta "github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"

	loader "github.com/cen-ngc5139/BeePF/loader/lib/src/cli"
	"go.uber.org/zap"
)

//go:generate sh -c "echo Generating for $TARGET_GOARCH"
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target $TARGET_GOARCH -go-package binary -output-dir ./binary -cc clang -no-strip kprobe_precpu ./bpf/kprobe_precpu.c -- -I../headers -Wno-address-of-packed-member

func main() {
	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("初始化日志失败: " + err.Error())
	}
	defer logger.Sync()

	config := &loader.Config{
		ObjectPath:    "./binary/kprobe_precpu_x86_bpfel.o",
		Logger:        logger,
		StructName:    "event",
		PollTimeout:   100 * time.Millisecond,
		IsEnableStats: true,
		StatsInterval: 1 * time.Second,
		ProgProperties: &meta.ProgProperties{
			CGroupPath: "/sys/fs/cgroup/unified",
		},
		// 设置用户自定义的 map 数据导出处理器
		UserExporterHandler: &export.MyCustomHandler{
			Logger: logger,
		},
		UserMetricsHandler: &metrics.DefaultHandler{
			Logger: logger,
		},
	}

	bpfLoader := loader.NewBPFLoader(config)

	err = bpfLoader.Init()
	if err != nil {
		logger.Fatal("初始化 BPF 加载器失败", zap.Error(err))
		return
	}

	err = bpfLoader.Load()
	if err != nil {
		logger.Fatal("加载 BPF 程序失败", zap.Error(err))
		return
	}

	if err := bpfLoader.Start(); err != nil {
		logger.Fatal("启动失败", zap.Error(err))
	}

	if err := bpfLoader.Stats(); err != nil {
		logger.Fatal("启动统计收集器失败", zap.Error(err))
	}

	if err := bpfLoader.Metrics(); err != nil {
		logger.Fatal("启动指标失败", zap.Error(err))
	}

	// 等待退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("正常关闭")
}
