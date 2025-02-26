package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"

	loader "github.com/cen-ngc5139/BeePF/loader/lib/src/cli"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"
	"go.uber.org/zap"
)

//go:generate sh -c "echo Generating for $TARGET_GOARCH"
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type span_info -target $TARGET_GOARCH -go-package binary -output-dir ./binary -cc clang -no-strip cgroup_skb ./bpf/cgroup_skb.c -- -I../headers -Wno-address-of-packed-member

func main() {
	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("初始化日志失败", zap.Error(err))
		return
	}
	defer logger.Sync()

	cgroupPath, err := detectCgroupPath()
	if err != nil {
		logger.Fatal("检测 cgroup 路径失败", zap.Error(err))
		return
	}

	config := &loader.Config{
		ObjectPath:    "./binary/cgroup_skb_x86_bpfel.o",
		Logger:        logger,
		StructName:    "pkt_count",
		PollTimeout:   100 * time.Millisecond,
		IsEnableStats: true,
		StatsInterval: 1 * time.Second,
		ProgProperties: &meta.ProgProperties{
			CGroupPath: cgroupPath,
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

// detectCgroupPath returns the first-found mount point of type cgroup2
// and stores it in the cgroupPath global variable.
func detectCgroupPath() (string, error) {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// example fields: cgroup2 /sys/fs/cgroup/unified cgroup2 rw,nosuid,nodev,noexec,relatime 0 0
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) >= 3 && fields[2] == "cgroup2" {
			return fields[1], nil
		}
	}

	return "", errors.New("cgroup2 not mounted")
}
