//go:build linux

// This program demonstrates attaching a fentry eBPF program to
// tcp_connect. It prints the command/IPs/ports information
// once the host sent a TCP SYN packet to a destination.
// It supports IPv4 at this example.
//
// Sample output:
//
// examples# go run -exec sudo ./fentry
// 2021/11/06 17:51:15 Comm   Src addr      Port   -> Dest addr        Port
// 2021/11/06 17:51:25 wget   10.0.2.15     49850  -> 142.250.72.228   443
// 2021/11/06 17:51:46 ssh    10.0.2.15     58854  -> 10.0.2.1         22
// 2021/11/06 18:13:15 curl   10.0.2.15     54268  -> 104.21.1.217     80
package main

import (
	"time"

	loader "github.com/cen-ngc5139/BeePF/loader/lib/src/cli"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"go.uber.org/zap"
)

//go:generate sh -c "echo Generating for $TARGET_GOARCH"
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type event -target $TARGET_GOARCH -go-package binary -output-dir ./binary -cc clang -no-strip fentry ./bpf/fentry.c -- -I../headers -Wno-address-of-packed-member

func main() {
	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("初始化日志失败", zap.Error(err))
		return
	}
	defer logger.Sync()

	config := &loader.Config{
		ObjectPath:  "./binary/fentry_x86_bpfel.o",
		Logger:      logger,
		StructName:  "event",
		PollTimeout: 100 * time.Millisecond,
		Properties: meta.Properties{
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
		logger.Fatal("start failed", zap.Error(err))
	}

	if err := bpfLoader.Stats(); err != nil {
		logger.Fatal("start stats collector failed", zap.Error(err))
	}

	if err := bpfLoader.Metrics(); err != nil {
		logger.Fatal("start metrics failed", zap.Error(err))
	}

	// 等待退出信号
	<-bpfLoader.Done()
	logger.Info("clean shutdown")
}
