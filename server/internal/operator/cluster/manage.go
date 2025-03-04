package cluster

import (
	k8scluster "github.com/cen-ngc5139/BeePF/server/internal/store/cluster"
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/cen-ngc5139/BeePF/server/pkg/utils"
)

type Cluster interface {
	Create() (err error)
	Update(id int) (err error)
	Get(id int) (cluster *models.Cluster, err error)
	List() (total int64, clusters []*models.Cluster, err error)
	Delete(id int) (err error)
	Bound() (err error)
	UnBound() (err error)
}

type Operator struct {
	QueryParma   *utils.Query
	Cluster      *models.Cluster
	ClusterStore *k8scluster.Store
	User         string
}

func NewOperator(user string) *Operator {
	return &Operator{
		User:         user,
		ClusterStore: &k8scluster.Store{},
	}
}

func (o *Operator) WithCluster(c *models.Cluster) *Operator {
	o.Cluster = c
	return o
}

func (o *Operator) WithQueryParma(q *utils.Query) *Operator {
	o.QueryParma = q
	return o
}
