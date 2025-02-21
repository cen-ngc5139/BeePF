package main

import (
	"fmt"
	_ "net/http/pprof"
	"time"

	"github.com/cen-ngc5139/BeePF/example/sched_wakeup/loader"
	"go.uber.org/zap"
)

//go:generate sh -c "echo Generating for $TARGET_GOARCH"
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type sched_latency_t -target $TARGET_GOARCH -go-package binary -output-dir ./binary -cc clang -no-strip Shepherd ./bpf/trace.c -- -I../headers -Wno-address-of-packed-member

func main() {
	fmt.Println("start")

	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		return
	}
	defer logger.Sync()

	config := &loader.Config{
		ObjectPath:  "./binary/shepherd_x86_bpfel.o",
		Logger:      logger,
		StructName:  "sched_latency_t",
		PollTimeout: 100 * time.Millisecond,
	}

	bpfLoader := loader.NewBPFLoader(config)

	err = bpfLoader.Init()
	if err != nil {
		fmt.Printf("初始化 BPF 加载器失败: %v\n", err)
		return
	}

	err = bpfLoader.Load()
	if err != nil {
		fmt.Printf("加载 BPF 程序失败: %v\n", err)
		return
	}

	if err := bpfLoader.Start(); err != nil {
		logger.Fatal("start failed", zap.Error(err))
	}

	// 等待退出信号，带超时
	select {
	case <-bpfLoader.Done():
		logger.Info("clean shutdown")
	case <-time.After(10 * time.Second):
		logger.Error("shutdown timeout")
	}
}
