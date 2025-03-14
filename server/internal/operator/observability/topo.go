package observability

import (
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	lib "github.com/cen-ngc5139/BeePF/loader/lib/src/observability/topology"
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
