package topology

import (
	"errors"
	"fmt"

	"github.com/cilium/ebpf"
	"golang.org/x/sys/unix"
)

// ListAllPrograms 通过 cilium/ebpf 获取当前节点上所有 prog 实例
func ListAllPrograms() (map[ebpf.ProgramID]*ebpf.ProgramInfo, error) {
	programs := make(map[ebpf.ProgramID]*ebpf.ProgramInfo)
	var nextID ebpf.ProgramID

	// 从 ID 0 开始遍历所有程序
	for {
		// 获取下一个程序 ID
		id, err := ebpf.ProgramGetNextID(nextID)
		if err != nil {
			// 如果是 ENOENT 错误，表示已经遍历完所有程序
			if errors.Is(err, unix.ENOENT) {
				break
			}
			return nil, fmt.Errorf("获取下一个程序 ID 失败: %w", err)
		}

		// 确保 ID 有变化，避免无限循环
		if id <= nextID {
			break
		}

		info, err := GetProgInfo(ebpf.ProgramID(id))
		if err != nil {
			// 如果无法获取程序信息，关闭程序并继续遍历
			nextID = id
			continue
		}

		// 添加程序名称到结果列表
		programs[id] = info

		// 更新下一个 ID
		nextID = id
	}

	return programs, nil
}

// ListAllMaps 列出所有 map
func ListAllMaps() (map[ebpf.MapID]*ebpf.MapInfo, error) {
	maps := make(map[ebpf.MapID]*ebpf.MapInfo)
	var nextID ebpf.MapID

	for {
		id, err := ebpf.MapGetNextID(nextID)
		if err != nil {
			if errors.Is(err, unix.ENOENT) {
				break
			}
			return nil, fmt.Errorf("获取下一个 map ID 失败: %w", err)
		}

		info, err := GetMapInfo(ebpf.MapID(id))
		if err != nil {
			nextID = id
			continue
		}

		maps[id] = info

		nextID = id
	}

	return maps, nil
}
