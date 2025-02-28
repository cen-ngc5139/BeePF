package loader

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"
	"go.uber.org/zap"
)

func TestBPFLoader_Init(t *testing.T) {

	cgroupPath, err := detectCgroupPath()
	if err != nil {
		fmt.Println("检测 cgroup 路径失败", zap.Error(err))
		return
	}

	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("初始化日志失败", zap.Error(err))
		return
	}
	defer logger.Sync()

	type fields struct {
		Config *Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "sched_wakeup",
			fields: fields{
				Config: &Config{
					ObjectPath:    "../../../../example/sched_wakeup/binary/shepherd_x86_bpfel.o",
					Logger:        logger,
					StructName:    "sched_latency_t",
					PollTimeout:   100 * time.Millisecond,
					IsEnableStats: true,
					StatsInterval: 1 * time.Second,
					// 设置用户自定义的 map 数据导出处理器
					UserExporterHandler: &export.MyCustomHandler{
						Logger: logger,
					},
					ProgProperties: &meta.ProgProperties{
						CGroupPath: cgroupPath,
					},
					// 设置用户自定义的 stats 数据导出处理器
					UserMetricsHandler: &metrics.DefaultHandler{
						Logger: logger,
					},
				},
			},
		},

		//{
		//	name: "cgroup_skb",
		//	fields: fields{
		//		Config: &Config{
		//			ObjectPath:    "../../../../example/cgroup_skb/binary/cgroup_skb_x86_bpfel.o",
		//			Logger:        logger,
		//			StructName:    "cgroup_skb_t",
		//			PollTimeout:   100 * time.Millisecond,
		//			IsEnableStats: true,
		//			StatsInterval: 1 * time.Second,
		//			ProgProperties: &meta.ProgProperties{
		//				CGroupPath: cgroupPath,
		//			},
		//			UserExporterHandler: &export.MyCustomHandler{
		//				Logger: logger,
		//			},
		//			UserMetricsHandler: &metrics.DefaultHandler{
		//				Logger: logger,
		//			},
		//		},
		//	},
		//},
		//
		//{
		//	name: "fentry",
		//	fields: fields{
		//		Config: &Config{
		//			ObjectPath:    "../../../../example/fentry/binary/fentry_x86_bpfel.o",
		//			Logger:        logger,
		//			StructName:    "event",
		//			PollTimeout:   100 * time.Millisecond,
		//			IsEnableStats: true,
		//			StatsInterval: 1 * time.Second,
		//			ProgProperties: &meta.ProgProperties{
		//				CGroupPath: cgroupPath,
		//			},
		//			UserExporterHandler: &export.MyCustomHandler{
		//				Logger: logger,
		//			},
		//			UserMetricsHandler: &metrics.DefaultHandler{
		//				Logger: logger,
		//			},
		//		},
		//	},
		//},
		//
		//{
		//	name: "kprobe",
		//	fields: fields{
		//		Config: &Config{
		//			ObjectPath:    "../../../../example/kprobe/binary/kprobe_x86_bpfel.o",
		//			Logger:        logger,
		//			PollTimeout:   100 * time.Millisecond,
		//			IsEnableStats: true,
		//			StatsInterval: 1 * time.Second,
		//			ProgProperties: &meta.ProgProperties{
		//				CGroupPath: cgroupPath,
		//			},
		//			UserExporterHandler: &export.MyCustomHandler{
		//				Logger: logger,
		//			},
		//			UserMetricsHandler: &metrics.DefaultHandler{
		//				Logger: logger,
		//			},
		//		},
		//	},
		//},
		//{
		//	name: "kprobe_precpu",
		//	fields: fields{
		//		Config: &Config{
		//			ObjectPath:    "../../../../example/kprobe_precpu/binary/kprobe_precpu_x86_bpfel.o",
		//			Logger:        logger,
		//			StructName:    "event",
		//			PollTimeout:   100 * time.Millisecond,
		//			IsEnableStats: true,
		//			StatsInterval: 1 * time.Second,
		//			ProgProperties: &meta.ProgProperties{
		//				CGroupPath: cgroupPath,
		//			},
		//			UserExporterHandler: &export.MyCustomHandler{
		//				Logger: logger,
		//			},
		//			UserMetricsHandler: &metrics.DefaultHandler{
		//				Logger: logger,
		//			},
		//		},
		//	},
		//},
		//
		//{
		//	name: "pin_path",
		//	fields: fields{
		//		Config: &Config{
		//			ObjectPath:    "../../../../example/kprobe_pin/binary/kprobepin_x86_bpfel.o",
		//			Logger:        logger,
		//			StructName:    "event",
		//			PollTimeout:   100 * time.Millisecond,
		//			IsEnableStats: true,
		//			StatsInterval: 1 * time.Second,
		//			ProgProperties: &meta.ProgProperties{
		//				PinPath:    "/sys/fs/bpf/kprobepin",
		//				CGroupPath: cgroupPath,
		//			},
		//			UserExporterHandler: &export.MyCustomHandler{
		//				Logger: logger,
		//			},
		//			UserMetricsHandler: &metrics.DefaultHandler{
		//				Logger: logger,
		//			},
		//		},
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpfLoader := NewBPFLoader(tt.fields.Config)

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
				logger.Fatal("export stats collector failed", zap.Error(err))
			}

			// 等待退出信号
			<-bpfLoader.Done()
			logger.Info("clean shutdown")
		})
	}
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
