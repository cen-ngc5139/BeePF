package topology

import (
	"errors"
	"fmt"

	"github.com/cilium/ebpf"
	"golang.org/x/sys/unix"
)

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

		m, err := ebpf.NewMapFromID(id)
		if err != nil {
			return nil, fmt.Errorf("获取 map 失败: %w", err)
		}

		info, err := m.Info()
		if err != nil {
			return nil, fmt.Errorf("获取 map 信息失败: %w", err)
		}

		maps[id] = info

		m.Close()

		nextID = id
	}

	return maps, nil
}
