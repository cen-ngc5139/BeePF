package cluster

import (
	"github.com/pkg/errors"
)

func (o *Operator) checkCluster() (err error) {
	basic := o.Cluster.ClusterBasic
	// 必填校验
	if err = basic.Validate(); err != nil {
		err = errors.Wrap(err, "真实集群校验失败")
		return
	}

	return
}
