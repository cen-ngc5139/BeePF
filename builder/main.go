package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// 模板定义
var templates = map[string]string{
	"Makefile": `GO := go
GO_BUILD = CGO_ENABLED=1 $(GO) build
GO_GENERATE = $(GO) generate
GO_TAGS ?=
TARGET_GOARCH ?= amd64,arm64
GOARCH ?= amd64
GOOS ?= linux
VERSION=$(shell git describe --tags --always)
# For compiling libpcap and CGO
CC ?= gcc


elf:
	TARGET_GOARCH=$(TARGET_GOARCH) $(GO_GENERATE)
    	CC=$(CC) GOARCH=$(TARGET_GOARCH) $(GO_BUILD) $(if $(GO_TAGS),-tags $(GO_TAGS)) \
    		-ldflags "-w -s "

build: elf
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_LDFLAGS='-g -lcapstone -static'   go build -tags=netgo,osusergo -gcflags "all=-N -l" -v  -o {{.ProgramName}}

dlv: build
	dlv --headless --listen=:2345 --api-version=2 exec ./{{.ProgramName}}	

run:  build
	./{{.ProgramName}}
`,

	"main.go": `package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"
	meta "github.com/cen-ngc5139/BeePF/loader/lib/src/meta"

	loader "github.com/cen-ngc5139/BeePF/loader/lib/src/cli"
	"go.uber.org/zap"
)

//go:generate sh -c "echo Generating for $TARGET_GOARCH"
{{if .StructName}}//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type {{.StructName}} -target $TARGET_GOARCH -go-package binary -output-dir ./binary -cc clang -no-strip {{.BinaryName}} ./bpf/{{.ProgramName}}.c -- -I../headers -Wno-address-of-packed-member{{else}}//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target $TARGET_GOARCH -go-package binary -output-dir ./binary -cc clang -no-strip {{.BinaryName}} ./bpf/{{.ProgramName}}.c -- -I../headers -Wno-address-of-packed-member{{end}}

func main() {
	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("初始化日志失败: " + err.Error())
	}
	defer logger.Sync()

	config := &loader.Config{
		ObjectPath:    "./binary/{{.BinaryName}}_x86_bpfel.o",
		Logger:        logger,
		StructName:    "{{.StructName}}",
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
`,

	"bpf_program.c": `// go:build ignore

#include "vmlinux.h"
#include "vmlinux-x86.h"
#include "bpf/bpf_helpers.h"

char __license[] SEC("license") = "Dual MIT/GPL";

// 定义 Map
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, u32);
    __type(value, u64);
} {{.MapName}} SEC(".maps");

// 定义事件结构体
struct {{.StructName}} {
    __u64 timestamp;      // 时间戳
    __u32 pid;            // 进程 ID
    __u32 tid;            // 线程 ID
    char comm[16];        // 进程名称
    __u64 counter;        // 计数器
} __attribute__((packed));

// 确保结构体类型信息被保留在 BTF 中
struct {{.StructName}} *unused_{{.StructNameLower}} __attribute__((unused));

struct
{
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, u32);
    __type(value, u64);
} kprobe_map SEC(".maps");

{{if eq .ProgramType "kprobe"}}
SEC("kprobe/rpc_exit_task")
int {{.HandlerName}}(struct pt_regs *regs)
{
        u32 key = 0;
    u64 initval = 1, *valp;

    valp = bpf_map_lookup_elem(&kprobe_map, &key);
    if (!valp)
    {
        bpf_map_update_elem(&kprobe_map, &key, &initval, BPF_ANY);
        return 0;
    }
    __sync_fetch_and_add(valp, 1);

    return 0;
}
{{end}}
`,
}

// 项目配置
type ProjectConfig struct {
	ProjectName     string
	ProgramName     string
	BinaryName      string
	StructName      string
	StructNameLower string
	MapName         string
	HandlerName     string
	ProgramType     string
}

func main() {
	// 解析命令行参数
	projectName := flag.String("name", "my-ebpf-project", "项目名称")
	programType := flag.String("type", "kprobe", "程序类型 (kprobe)")
	structName := flag.String("struct", "event", "事件结构体名称 (可选)")
	flag.Parse()

	// 验证程序类型
	validTypes := map[string]bool{
		"kprobe": true,
	}
	if !validTypes[*programType] {
		fmt.Printf("错误: 不支持的程序类型 '%s'。支持的类型: kprobe\n", *programType)
		os.Exit(1)
	}

	// 创建项目配置
	config := ProjectConfig{
		ProjectName:     *projectName,
		ProgramName:     strings.ToLower(*projectName),
		BinaryName:      strings.ToLower(*projectName),
		StructName:      *structName,
		StructNameLower: strings.ToLower(*structName),
		MapName:         "data_map",
		HandlerName:     fmt.Sprintf("handle_%s", strings.ToLower(*programType)),
		ProgramType:     *programType,
	}

	// 创建项目目录结构
	dirs := []string{
		*projectName,
		filepath.Join(*projectName, "bpf"),
		filepath.Join(*projectName, "binary"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("创建目录失败 %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// 生成文件
	files := map[string]string{
		filepath.Join(*projectName, "Makefile"):                     "Makefile",
		filepath.Join(*projectName, "main.go"):                      "main.go",
		filepath.Join(*projectName, "bpf", config.ProgramName+".c"): "bpf_program.c",
	}

	for path, templateName := range files {
		if err := generateFile(path, templateName, config); err != nil {
			fmt.Printf("生成文件失败 %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	fmt.Printf("项目 '%s' 已成功创建!\n", *projectName)
	fmt.Printf("使用以下命令构建和运行:\n")
	fmt.Printf("  cd %s\n", *projectName)
	fmt.Printf("  make build\n")
	fmt.Printf("  sudo ./run\n")
}

// 生成文件
func generateFile(path, templateName string, config ProjectConfig) error {
	tmpl, err := template.New(templateName).Parse(templates[templateName])
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, config)
}
