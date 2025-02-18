# BeePF

BeePF 是一个用 Go 语言编写的 eBPF 程序加载器和运行时框架。它提供了一套完整的工具链，用于加载、管理和监控 eBPF 程序。

## 特性

- 支持自动加载和管理 BTF (BPF Type Format) 信息
- 提供多种数据导出格式 (JSON、纯文本、原始数据、直方图等)
- 灵活的事件处理机制
- 支持 Map 数据采样和导出
- 内置性能监控和调试功能

## 安装

确保你的系统满足以下要求：

- Go 1.16 或更高版本
- Linux 内核 5.4 或更高版本（支持 BTF）
- `clang` 和 `llvm` 用于编译 eBPF 程序

通过 Go 工具链安装：

```bash
go get github.com/cen-ngc5139/BeePF
```

## 快速开始

1. 创建一个简单的 eBPF 程序：

```c
// example.bpf.c
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

struct event {
    __u32 pid;
    __u32 uid;
    char comm[16];
};

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(int));
    __uint(value_size, sizeof(int));
} events SEC(".maps");

SEC("kprobe/do_sys_open")
int kprobe_do_sys_open(struct pt_regs *ctx)
{
    struct event e = {};
    e.pid = bpf_get_current_pid_tgid() >> 32;
    bpf_get_current_comm(&e.comm, sizeof(e.comm));
    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &e, sizeof(e));
    return 0;
}
```

2. 使用 BeePF 加载和运行程序：

```go
package main

import (
    "github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
    "github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton"
)

func main() {
    // 加载 eBPF 对象文件
    pkg, err := meta.GenerateComposedObject("example.bpf.o")
    if err != nil {
        panic(err)
    }

    // 创建骨架构建器
    builder := skeleton.FromJsonPackage(pkg, "")
    
    // 构建和加载 eBPF 程序
    skel, err := builder.Build()
    if err != nil {
        panic(err)
    }

    // 运行程序
    if err := skel.Load(); err != nil {
        panic(err)
    }
}
```

## 架构

BeePF 主要由以下组件组成：

- **Loader**: eBPF 程序加载器
- **Skeleton**: 程序骨架生成器和管理器
- **BTF**: BTF 信息处理器
- **Export**: 数据导出和处理模块
- **Container**: ELF 和 BTF 容器管理

## 配置

BeePF 支持通过环境变量和配置文件进行配置：

- `BTF_FILE_PATH`: 指定 BTF 文件路径
- `VMLINUX_BTF_PATH`: 系统 BTF 文件路径（默认：/sys/kernel/btf/vmlinux）

## 贡献

欢迎提交 Pull Request 和 Issue！在提交代码前，请确保：

1. 代码通过所有测试
2. 新功能包含相应的测试用例
3. 更新相关文档

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。