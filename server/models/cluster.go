package models

import (
	"time"

	"github.com/pkg/errors"
)

type EnvironmentType string

// 测试，开发，验收，生产
const (
	TestEnvironment   EnvironmentType = "test"
	DeployEnvironment EnvironmentType = "dev"
	SitEnvironment    EnvironmentType = "sit"
	ProdEnvironment   EnvironmentType = "prod"
)

type Status int

const (
	ClusterUp   Status = 1
	ClusterDown Status = 0
)

type ClusterBasic struct {
	Name        string          `gorm:"column:cluster_name" json:"name"`
	CnName      string          `gorm:"column:cn_name" json:"cnname"`
	Master      string          `gorm:"column:cluster_master" json:"master"`
	KubeConfig  string          `gorm:"column:kube_config" json:"kubeconfig"`
	Status      Status          `gorm:"column:cluster_status" json:"status"` // status 1 is up 0 is down
	Desc        string          `gorm:"column:cluster_desc" json:"desc"`
	Creator     string          ` json:"creator"`
	Environment EnvironmentType `json:"environment"`
}

type Cluster struct {
	ClusterBasic
	Id        int       `json:"id"`
	Deleted   bool      `gorm:"default false" json:"deleted"`
	CreatedAt time.Time `gorm:"column:created_time" json:"createdat"`
	UpdateAt  time.Time `gorm:"column:last_update_time" json:"updateat"`
}

const (
	ClusterTable = "cluster"

	ClusterMaster     = "cluster_master"
	ClusterKubeConfig = "kube_config"
	ClusterDesc       = "cluster_desc"
	ClusterStatus     = "cluster_status"
	ClusterEnv        = "environment"
	ClusterDeleted    = "deleted"
	ClusterUpdateAt   = "last_update_time"
)

func (m *Cluster) TableName() string {
	return "cluster"
}

func (m *Cluster) WitchCreator(user string) *Cluster {
	m.Creator = user
	m.UpdateAt = time.Now()
	return m
}

func (i *ClusterBasic) Validate() error {
	if len(i.Name) == 0 || len(i.CnName) == 0 || len(i.Master) == 0 ||
		len(i.Environment) == 0 {
		return errors.New(" 集群信息填写有误，请校验！")
	}

	return nil
}

func (m *Cluster) UpdateCluster(current *Cluster) *Cluster {
	current.Master = m.Master
	current.CnName = m.CnName
	current.KubeConfig = m.KubeConfig
	current.Desc = m.Desc
	current.Status = m.Status
	current.Environment = m.Environment
	current.UpdateAt = time.Now()

	return current
}
