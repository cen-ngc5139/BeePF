package meta

import "github.com/cilium/ebpf"

type Properties struct {
	// Maps 映射列表
	Maps map[string]Map

	// Programs 程序列表
	Programs map[string]Program
}

type Map struct {
	// Name 映射名称
	Name string

	// ExportHandler 导出处理器
	ExportHandler EventHandler
}

type Program struct {
	// Name 程序名称
	Name string

	// CGroupPath 用于 cgroup 程序的 cgroup 路径
	CGroupPath string

	// PinPath 用于指定 eBPF 程序的 pin 路径，下次加载时从该路径加载
	PinPath string

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
