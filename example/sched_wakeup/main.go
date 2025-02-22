package main

import (
	"time"

	loader "github.com/cen-ngc5139/BeePF/loader/lib/src/cli"
	"go.uber.org/zap"
)

//go:generate sh -c "echo Generating for $TARGET_GOARCH"
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type sched_latency_t -target $TARGET_GOARCH -go-package binary -output-dir ./binary -cc clang -no-strip Shepherd ./bpf/trace.c -- -I../headers -Wno-address-of-packed-member

func main() {
	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("初始化日志失败", zap.Error(err))
		return
	}
	defer logger.Sync()

	config := &loader.Config{
		ObjectPath:    "./binary/shepherd_x86_bpfel.o",
		Logger:        logger,
		StructName:    "sched_latency_t",
		PollTimeout:   100 * time.Millisecond,
		IsEnableStats: true,
		StatsInterval: 1 * time.Second,
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
		logger.Fatal("start failed", zap.Error(err))
	}

	if err := bpfLoader.StatsCollector.Start(); err != nil {
		logger.Fatal("start stats collector failed", zap.Error(err))
	}

	// 定时从 stats collector 中获取 stats 信息
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			programs, err := bpfLoader.StatsCollector.GetPrograms()
			if err != nil {
				logger.Error("获取 stats 信息失败", zap.Error(err))
			}
			for _, program := range programs {
				logger.Info("program", zap.Any("program", program))
			}
		}
	}()

	// 等待退出信号
	<-bpfLoader.Done()
	logger.Info("clean shutdown")
}
