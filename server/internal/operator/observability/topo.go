package observability

import (
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	lib "github.com/cen-ngc5139/BeePF/loader/lib/src/observability/topology"
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/pkg/errors"
)

type Topo struct{}

func NewTopo() *Topo {
	return &Topo{}
}

func (t *Topo) GetTopo() (meta.Topology, error) {
	topology, err := lib.MergeTopology()
	if err != nil {
		return topology, errors.Wrap(err, "获取程序失败")
	}

	return topology, nil
}

func (t *Topo) ListProgs() ([]models.ProgramInfoWrapper, error) {
	progs, err := lib.ListAllPrograms()
	if err != nil {
		return nil, errors.Wrap(err, "获取程序失败")
	}

	programInfos := make([]models.ProgramInfoWrapper, 0, len(progs))
	for id, info := range progs {
		wrapper := models.ProgramInfoWrapper{
			ID:   id,
			Name: info.Name,
			Type: info.Type,
			Tag:  info.Tag,
		}

		maps, ok := info.MapIDs()
		if ok {
			wrapper.Maps = maps
		}

		btfID, ok := info.BTFID()
		if ok {
			wrapper.BTF = btfID
		}

		loadTime, ok := info.LoadTime()
		if ok {
			wrapper.LoadTime = loadTime
		}

		createdByUID, haveCreatedByUID := info.CreatedByUID()
		if haveCreatedByUID {
			wrapper.CreatedByUID = createdByUID
		}

		programInfos = append(programInfos, wrapper)
	}

	return programInfos, nil
}
