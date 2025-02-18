package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cilium/ebpf/perf"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"
	"github.com/cilium/ebpf"
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

	// 生成组合对象
	fmt.Println("正在生成组合对象...")
	pkg, err := meta.GenerateComposedObject("./binary/shepherd_x86_bpfel.o")
	if err != nil {
		fmt.Printf("生成组合对象失败: %v\n", err)
		return
	}
	fmt.Println("组合对象生成成功")

	// 构建预加载骨架
	fmt.Println("正在构建预加载骨架...")
	preLoadBpfSkeleton, err := skeleton.FromJsonPackage(pkg, "./binary").Build()
	if err != nil {
		fmt.Printf("构建预加载骨架失败: %v\n", err)
		return
	}
	fmt.Println("预加载骨架构建成功")

	// 加载并附加 eBPF 程序
	fmt.Println("正在加载并附加 eBPF 程序...")
	skel, err := preLoadBpfSkeleton.LoadAndAttach()
	if err != nil {
		fmt.Printf("加载并附加 eBPF 程序失败: %v\n", err)
		return
	}
	fmt.Println("eBPF 程序加载并附加成功")
	defer skel.Collection.Close()

	// 设置环形缓冲区
	for _, m := range skel.Collection.Maps {
		if m.Type() != ebpf.PerfEventArray {
			continue
		}

		// 创建环形缓冲区读取器
		perfReader, err := perf.NewReader(m, os.Getpagesize())
		if err != nil {
			log.Fatalf("创建环形缓冲区读取器失败: %v", err)
		}
		defer perfReader.Close()

		// 设置事件导出器
		ee := export.NewEventExporterBuilder().
			SetExportFormat(export.FormatJson).
			SetUserContext(export.NewUserContext(0)).
			SetEventHandler(&export.MyCustomHandler{Logger: logger})

		// 查找并处理调度延迟事件
		for _, v := range skel.Collection.Variables {
			structType, err := skeleton.FindStructType(v.Type())
			if err != nil {
				log.Printf("查找结构体类型失败: %v", err)
				continue
			}

			if structType.Name != "sched_latency_t" {
				continue
			}

			// 构建事件导出器
			exporter, err := ee.BuildForSingleValueWithTypeDescriptor(
				&export.BTFTypeDescriptor{
					Type: structType,
					Name: structType.TypeName(),
				},
				skel.Btf,
			)
			if err != nil {
				log.Fatalf("构建事件导出器失败: %v", err)
			}

			// 创建 JSON 处理器
			jsonHandler := export.NewJsonExportEventHandler(exporter)

			// 创建环形缓冲区轮询器
			p := &skeleton.PerfEventPoller{
				Reader:    perfReader,
				Processor: jsonHandler,
				Timeout:   100 * time.Millisecond,
			}

			// 创建程序轮询器
			pp := skeleton.NewProgramPoller(100 * time.Millisecond)
			pp.StartPolling("sched_wakeup", p.GetPollFunc(), func(err error) {
				log.Printf("轮询错误: %v", err)
			})
			defer pp.Stop()

			// 等待中断信号
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan

			log.Println("收到中断信号，正在退出...")
			return
		}
	}
}
