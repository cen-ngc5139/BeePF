package topology

import (
	"fmt"

	"github.com/cilium/ebpf"
)

// GetProgInfo 获取程序信息
func GetProgInfo(progID ebpf.ProgramID) (*ebpf.ProgramInfo, error) {
	// 根据 ID 获取程序信息
	prog, err := ebpf.NewProgramFromID(progID)
	if err != nil {
		return nil, err
	}

	info, err := prog.Info()
	if err != nil {
		return nil, err
	}

	fullName, err := GetFullName(info)
	if err == nil {
		info.Name = fullName
	}

	prog.Close()

	return info, nil
}

// GetMapInfo 获取 map 信息
func GetMapInfo(mapID ebpf.MapID) (*ebpf.MapInfo, error) {
	m, err := ebpf.NewMapFromID(mapID)
	if err != nil {
		return nil, err
	}

	info, err := m.Info()
	if err != nil {
		return nil, err
	}

	m.Close()

	return info, nil
}

// GetFullName 返回程序的完整名称。
// 如果程序有关联的 BTF ID 和函数信息，将尝试从 BTF 中获取函数名称。
// 否则，返回程序加载时设置的名称。
//
// 这类似于 bpftool 中的 get_prog_full_name 函数。
func GetFullName(pi *ebpf.ProgramInfo) (string, error) {
	funcInfos, err := pi.FuncInfos()
	if err != nil {
		return pi.Name, fmt.Errorf("get func info: %w", err)
	}

	// 如果没有函数，返回当前名称
	if len(funcInfos) == 0 {
		return pi.Name, nil
	}

	// 获取第一个函数的 BTF 类型
	funcType := funcInfos[0].Func
	if funcType == nil {
		return pi.Name, nil
	}

	// 获取函数名称
	funcName := funcType.Name
	if funcName == "" {
		return pi.Name, nil
	}

	return funcName, nil
}
