package metrics

import (
	"fmt"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/observability/topology"
	"go.uber.org/zap"
)

type NodeMetricsCollector struct {
	*topology.NodeMetricsCollector
}

func NewNodeMetricsCollector(interval time.Duration, logger *zap.Logger) (*NodeMetricsCollector, error) {
	collector, err := topology.NewNodeMetricsCollector(interval, logger)
	if err != nil {
		return nil, fmt.Errorf("fail to create node metrics collector: %v", err)
	}

	if err := collector.Start(); err != nil {
		return nil, fmt.Errorf("fail to start node metrics collector: %v", err)
	}

	return &NodeMetricsCollector{
		NodeMetricsCollector: collector,
	}, nil
}

func (c *NodeMetricsCollector) GetMetrics() (map[uint32]*meta.MetricsStats, error) {
	curr, err := c.NodeMetricsCollector.GetMetrics()
	if err != nil {
		return nil, fmt.Errorf("fail to get node metrics: %v", err)
	}

	for id, s := range curr {
		if s.CPUTimePercent == 0 {
			delete(curr, id)
		}
	}

	return curr, nil
}
