package topology

import (
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
