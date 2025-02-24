package loader

import (
	"testing"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"
	"go.uber.org/zap"
)

func TestBPFLoader_Init(t *testing.T) {
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
					// 设置用户自定义的 stats 数据导出处理器
					UserMetricsHandler: &metrics.DefaultHandler{
						Logger: logger,
					},
				},
			},
		},
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
