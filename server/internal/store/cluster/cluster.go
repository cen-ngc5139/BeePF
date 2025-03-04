package k8scluster

import (
	"fmt"
	"strings"

	"github.com/cen-ngc5139/BeePF/server/database"
	"github.com/cen-ngc5139/BeePF/server/models"

	"github.com/pkg/errors"
)

func (s *Store) Create(cluster *models.Cluster) (err error) {
	return database.DB.Table(models.ClusterTable).Create(cluster).Error
}

func (s *Store) List(pageSize, pageNum int, attachs map[string]interface{}) (total int64, clusters []*models.Cluster, err error) {
	pv := map[string]interface{}{"deleted": false}
	// merge key and value
	for key, v := range attachs {
		pv[key] = v
	}
	keys, values := make([]string, 0, len(pv)), make([]interface{}, 0, len(pv))
	for key, value := range pv {
		keys = append(keys, fmt.Sprintf("`%s` = ?", key))
		values = append(values, value)
	}
	db := database.DB.Table(models.ClusterTable).Where(strings.Join(keys, " AND "), values...)
	if _err := db.Count(&total).Error; _err != nil {
		err = errors.Wrap(_err, "获取真实集群总数失败")
	}
	if pageSize > 0 && pageNum > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&clusters).Error
	} else {
		err = db.Find(&clusters).Error
	}

	if err != nil {
		err = errors.Wrap(err, "读取数据库失败")
		return
	}

	return
}

func (s *Store) Get(id int) (cluster models.Cluster, err error) {
	err = database.DB.Table(models.ClusterTable).
		Where("deleted = ? AND id = ?", false, id).First(&cluster).Error
	if err != nil {
		err = errors.Wrap(err, "读取数据库失败")
		return
	}

	err = database.DB.Model(&cluster).Preload("Clusters").First(&cluster).Error
	if err != nil {
		err = errors.Wrap(err, "真实集群关联查询虚拟集群失败")
		return
	}

	return
}

func (s *Store) Count() (total int64, err error) {
	if err = database.DB.Table(models.ClusterTable).Where("deleted = ?", false).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "获取真实集群总数失败")
		return
	}

	return
}

func (s *Store) Update(cluster *models.Cluster) (err error) {
	return database.DB.Table(models.ClusterTable).Where("deleted = ? AND id = ?", false, cluster.Id).
		UpdateColumns(BuildUpdateObj(cluster)).Error
}

func (s *Store) Delete(id int) (err error) {
	return database.DB.Table(models.ClusterTable).Where("deleted = ? AND id = ?", false, id).
		UpdateColumns(BuildDeleteObj()).Error
}

func (s *Store) Find(parma models.ClusterBasic) (clusters []*models.Cluster, count int64, err error) {
	qs := database.DB.Table(models.ClusterTable).Where("deleted = ?", false)

	if len(parma.Name) != 0 {
		qs = qs.Where("cluster_name", parma.Name)
	}

	if err = qs.Count(&count).Error; err != nil {
		return
	}

	if err = qs.Find(&clusters).Error; err != nil {
		return
	}

	return
}

func BuildUpdateObj(vc *models.Cluster) (obj map[string]interface{}) {
	obj = make(map[string]interface{})
	obj[models.ClusterMaster] = vc.Master
	obj[models.ClusterKubeConfig] = vc.KubeConfig
	obj[models.ClusterDesc] = vc.Desc
	obj[models.ClusterStatus] = vc.Status
	obj[models.ClusterEnv] = vc.Environment
	obj[models.ClusterUpdateAt] = vc.UpdateAt

	return
}

func BuildDeleteObj() (obj map[string]interface{}) {
	obj = make(map[string]interface{})
	obj[models.ClusterDeleted] = true
	return
}
