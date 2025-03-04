package service

import (
	"errors"

	"github.com/cen-ngc5139/BeePF/server/internal/operator/cluster"
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/cen-ngc5139/BeePF/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Cluster struct{}

func (m *Cluster) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		basic := &models.ClusterBasic{}
		if err := c.BindJSON(basic); utils.HandleError(c, err) {
			return
		}

		cluster := cluster.NewOperator().
			WithCluster(&models.Cluster{ClusterBasic: *basic})
		if err := cluster.Create(); utils.HandleError(c, err) {
			return
		}

		data := &map[string]interface{}{}
		utils.HandleResult(c, data)
	}
}

func (m *Cluster) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		basic := &models.ClusterBasic{}
		if err := c.BindJSON(basic); utils.HandleError(c, err) {
			return
		}

		id := utils.GetParamIntItem("clusterId", c)
		if id == 0 {
			utils.ResponseErr(c, errors.New("真实集群编号不合法"))
			return
		}

		cluster := cluster.NewOperator().
			WithCluster(&models.Cluster{ClusterBasic: *basic})
		if err := cluster.Update(id); utils.HandleError(c, err) {
			return
		}

		data := &map[string]interface{}{}
		utils.HandleResult(c, data)
	}
}

func (m *Cluster) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := utils.GetParamIntItem("clusterId", c)

		if id == 0 {
			utils.ResponseErr(c, errors.New("真实集群编号不合法"))
			return
		}

		operator := cluster.NewOperator()
		if err := operator.Delete(id); utils.HandleError(c, err) {
			return
		}

		data := &map[string]interface{}{}
		utils.HandleResult(c, data)
	}
}

func (m *Cluster) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := utils.GetParamIntItem("clusterId", c)

		if id == 0 {
			utils.ResponseErr(c, errors.New("真实集群编号不合法"))
			return
		}

		operator := cluster.NewOperator()
		cluster, err := operator.Get(id)
		if utils.HandleError(c, err) {
			return
		}

		data := &map[string]interface{}{"cluster": cluster}
		utils.HandleResult(c, data)
	}
}

func (m *Cluster) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageSize, pageNum := utils.GetPageInfo(c)
		parma := utils.NewQueryParma(pageSize, pageNum)

		cluster := cluster.NewOperator().WithQueryParma(parma)
		total, clusters, err := cluster.List(nil)
		if utils.HandleError(c, err) {
			return
		}

		// 脱敏
		for i := range clusters {
			clusters[i].KubeConfig = ""
		}
		data := map[string]interface{}{"list": clusters, "total": total}
		utils.HandleResult(c, &data)
	}
}

func (m *Cluster) GetClustersByParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterName, clusterID := c.Query("clusterName"), utils.GetQueryIntItem("clusterId", "-1", c)
		pageSize, pageNum := utils.GetPageInfo(c)
		parma := utils.NewQueryParma(pageSize, pageNum)
		cluster := cluster.NewOperator().WithQueryParma(parma)
		attachs := map[string]interface{}{}
		if clusterName != "" {
			attachs["cluster_name"] = clusterName
		}

		if clusterID > 0 {
			attachs["id"] = clusterID
		}

		total, clusters, err := cluster.List(attachs)
		if utils.HandleError(c, err) {
			return
		}

		data := map[string]interface{}{"list": clusters, "total": total}
		utils.HandleResult(c, &data)
	}
}
