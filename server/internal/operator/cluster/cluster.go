package cluster

import (
	s "github.com/cen-ngc5139/BeePF/server/internal/store/cluster"
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/pkg/errors"
)

var store = &s.Store{}

func (o *Operator) Create() (err error) {
	if err = o.checkCluster(); err != nil {
		err = errors.Wrapf(err, "检查真实集群参数失败")
		return
	}

	err = o.ClusterStore.Create(o.Cluster.WitchCreator(o.User))
	if err != nil {
		err = errors.Wrapf(err, "新增真实集群 %s 失败", o.Cluster.Name)
		return
	}

	return
}

func (o *Operator) Update(id int) (err error) {
	if err = o.checkCluster(); err != nil {
		err = errors.Wrapf(err, "检查真实集群参数失败")
		return
	}

	var current models.Cluster
	current, err = o.ClusterStore.Get(id)
	if err != nil {
		err = errors.Wrapf(err, "获取虚拟集群信息失败")
		return
	}

	err = o.ClusterStore.Update(o.Cluster.UpdateCluster(&current))
	if err != nil {
		err = errors.Wrapf(err, "更新真实集群 %s 失败", current.Name)
		return
	}

	return
}

func (o *Operator) Delete(id int) (err error) {
	var current models.Cluster
	current, err = o.ClusterStore.Get(id)
	if err != nil {
		err = errors.Wrapf(err, "获取真实集群信息失败")
		return
	}

	err = o.ClusterStore.Delete(id)
	if err != nil {
		err = errors.Wrapf(err, "删除虚拟集群 %s 失败", current.Name)
		return
	}

	return
}

func (o *Operator) Get(id int) (cluster models.Cluster, err error) {
	return o.ClusterStore.Get(id)
}

func (o *Operator) List(attachs map[string]interface{}) (total int64, clusters []*models.Cluster, err error) {
	total, clusters, err = o.ClusterStore.List(o.QueryParma.PageSize, o.QueryParma.PageNum, attachs)
	if err != nil {
		err = errors.Wrap(err, "获取真实集群失败")
		return
	}

	if o.QueryParma.IsAdmin {
		return
	}

	for _, cluster := range clusters {
		isFound := false
		for _, pm := range o.QueryParma.Authorized {
			if pm == cluster.Name {
				isFound = true
				break
			}
		}

		if !isFound {
			cluster.KubeConfig = ""
		}
	}

	return
}

func (o *Operator) GetK8ClusterByName(name string) (cluster *models.Cluster, has bool, err error) {
	paramsMap := map[string]interface{}{"cluster_name": name}
	total, list, err := o.List(paramsMap)
	if err != nil {
		return
	}

	if total != 1 {
		err = errors.Errorf("真实集群 %s 查询结果不唯一，请联系管理员", name)
		return
	}

	for _, v := range list {
		if v.Name == name {
			has = true
			cluster = v
			break
		}
	}

	return
}
