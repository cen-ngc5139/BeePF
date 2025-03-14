package topology

import (
	"log"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/pkg/errors"
)

func MergeTopology() (meta.Topology, error) {
	topology := meta.NewTopology()
	progs, err := ListAllPrograms()
	if err != nil {
		return topology, errors.Wrap(err, "获取程序失败")
	}

	maps, err := ListAllMaps()
	if err != nil {
		return topology, errors.Wrap(err, "获取 map 失败")
	}

	for _, info := range progs {
		progID, ok := info.ID()
		if !ok {
			log.Fatalf("获取程序 ID 失败: %v", err)
			continue
		}

		mapIDs, ok := info.MapIDs()
		if !ok {
			log.Fatalf("获取 map ID 失败: %v", err)
			continue
		}

		for _, mapID := range mapIDs {
			topology.Edges.AddEdge(meta.TopologyEdge{
				ProgID: uint32(progID),
				MapID:  uint32(mapID),
			})
		}

		topology.ProgNodes.AddProgNode(meta.ProgNode{
			ID:   uint32(progID),
			Name: info.Name,
		})
	}

	for _, info := range maps {
		mapID, ok := info.ID()
		if !ok {
			log.Fatalf("获取 map ID 失败: %v", err)
			continue
		}

		topology.MapNodes.AddMapNode(meta.MapNode{
			ID:   uint32(mapID),
			Name: info.Name,
		})
	}

	return topology, nil
}
