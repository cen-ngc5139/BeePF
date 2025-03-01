package meta

import (
	"time"

	"github.com/cilium/ebpf"
)

// Properties 配置结构
// 事件处理器和指标处理器存在全局配置和 map 或程序的局部配置
// 局部配置会覆盖全局配置，如果 map 或程序没有配置，则使用全局配置
type Properties struct {
	// Maps 映射列表
	Maps map[string]*Map

	// Programs 程序列表
	Programs map[string]*Program

	// Stats 统计配置
	Stats *Stats

	// EventHandler 全局事件处理器
	EventHandler EventHandler

	// MetricsHandler 全局指标处理器
	MetricsHandler MetricsHandler
}

type Map struct {
	// Name 映射名称
	Name string

	// ExportHandler 导出处理器
	ExportHandler EventHandler

	// Properties 映射配置
	Properties *MapProperties
}

type MapProperties struct {
	// PinPath 用于指定 eBPF 映射的 pin 路径，下次加载时从该路径加载
	PinPath string
}

type Stats struct {
	// Interval 统计间隔
	Interval time.Duration

	// Handler 指标处理器
	Handler MetricsHandler
}

type Program struct {
	// Name 程序名称
	Name string

	// Properties 程序配置
	Properties *ProgramProperties
}

type ProgramProperties struct {
	// CGroupPath 用于 cgroup 程序的 cgroup 路径
	CGroupPath string

	// PinPath 用于指定 eBPF 程序的 pin 路径，下次加载时从该路径加载
	PinPath string

	// LinkPinPath 用于指定 eBPF 程序的 pin 路径，下次加载时从该路径加载
	LinkPinPath string

	// Tc 流量控制配置
	Tc *TCCLS
}

type TCCLS struct {
	// Ifindex 接口索引
	Ifindex int32

	// Ifname 接口名称
	Ifname string

	// AttachType 附加点类型
	AttachType ebpf.AttachType
}
