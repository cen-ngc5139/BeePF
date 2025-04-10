package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	loader "github.com/cen-ngc5139/BeePF/loader/lib/src/cli"
	"github.com/cilium/ebpf/rlimit"
	"go.uber.org/zap"
)

//go:generate sh -c "echo Generating for $TARGET_GOARCH"
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target $TARGET_GOARCH -go-package binary -output-dir ./binary -cc clang -no-strip tp_unlinkat ./bpf/tp_unlinkat.c -- -I../headers -Wno-address-of-packed-member

func main() {
	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("初始化日志失败: " + err.Error())
	}
	defer logger.Sync()

	// 移除 eBPF 程序的内存限制
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("Failed to remove memlock limit: %v\n", err)
		os.Exit(1)
	}

	config := &loader.Config{
		ObjectPath:  "./binary/tp_unlinkat_x86_bpfel.o",
		Logger:      logger,
		PollTimeout: 100 * time.Millisecond,
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
