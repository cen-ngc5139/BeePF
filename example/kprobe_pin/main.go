package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	loader "github.com/cen-ngc5139/BeePF/loader/lib/src/cli"
	meta "github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"go.uber.org/zap"
)

//go:generate sh -c "echo Generating for $TARGET_GOARCH"
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target $TARGET_GOARCH -go-package binary -output-dir ./binary -cc clang -no-strip kprobepin ./bpf/kprobepin.c -- -I../headers -Wno-address-of-packed-member

func main() {
	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("初始化日志失败: " + err.Error())
	}
	defer logger.Sync()

	config := &loader.Config{
		ObjectPath:  "./binary/kprobepin_x86_bpfel.o",
		Logger:      logger,
		PollTimeout: 100 * time.Millisecond,
		Properties: meta.Properties{
			Programs: map[string]*meta.Program{
				"rpc_exit_task": {
					Name:       "rpc_exit_task",
					Properties: &meta.ProgramProperties{PinPath: "/sys/fs/bpf/kprobepin/rpc_exit_task"},
				},
			},
			Maps: map[string]*meta.Map{
				"kprobe_map": {
					Name:       "kprobe_map",
					Properties: &meta.MapProperties{PinPath: "/sys/fs/bpf/kprobepin/"},
				},
			},
			Stats: &meta.Stats{
				Interval: 1 * time.Second,
				Handler:  metrics.NewDefaultHandler(logger),
			},
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
