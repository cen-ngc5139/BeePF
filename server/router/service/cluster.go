package service

import (
	"errors"

	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/cen-ngc5139/BeePF/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (m *Cluster) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		basic := &models.ClusterBasic{}
		if err := c.BindJSON(basic); utils.HandleError(c, err) {
			return
		}

		cluster := k8scluster.NewOperator(utils.PickUserInSession(c)).
			WithK8sCluster(&models.K8sCluster{K8sClusterBasic: *basic})
		if err := cluster.Create(); utils.HandleError(c, err) {
			return
		}

		data := &map[string]interface{}{}
		utils.HandleResult(c, data)
	}
}

func (m *K8sCluster) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		basic := &models.K8sClusterBasic{}
		if err := c.BindJSON(basic); utils.HandleError(c, err) {
			return
		}

		id := utils.GetParamIntItem("clusterId", c)
		if id == 0 {
			utils.ResponseErr(c, errors.New("真实集群编号不合法"))
			return
		}

		cluster := k8scluster.NewOperator(utils.PickUserInSession(c)).
			WithK8sCluster(&models.K8sCluster{K8sClusterBasic: *basic})
		if err := cluster.Update(id); utils.HandleError(c, err) {
			return
		}

		data := &map[string]interface{}{}
		utils.HandleResult(c, data)
	}
}

func (m *K8sCluster) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.PickUserInSession(c)
		id := utils.GetParamIntItem("clusterId", c)

		if id == 0 {
			utils.ResponseErr(c, errors.New("真实集群编号不合法"))
			return
		}

		operator := k8scluster.NewOperator(user)
		if err := operator.Delete(id); utils.HandleError(c, err) {
			return
		}

		data := &map[string]interface{}{}
		utils.HandleResult(c, data)
	}
}

func (m *K8sCluster) Get() gin.HandlerFunc {
	return func(c *gin.Context) {

		user := utils.PickUserInSession(c)
		id := utils.GetParamIntItem("clusterId", c)

		if id == 0 {
			utils.ResponseErr(c, errors.New("真实集群编号不合法"))
			return
		}

		operator := k8scluster.NewOperator(user)
		cluster, err := operator.Get(id)
		if utils.HandleError(c, err) {
			return
		}

		// 真实集群kubeconfig脱敏
		cluster.KubeConfig = ""
		data := &map[string]interface{}{"k8s_cluster": cluster}
		utils.HandleResult(c, data)
	}
}

func (m *K8sCluster) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageSize, pageNum := utils.GetPageInfo(c)
		user := utils.PickUserInSession(c)
		authed := utils.PickUserAuthorized(c)
		parma := utils.NewQueryParma(pageSize, pageNum, utils.GetIsAdmin(c), authed)

		cluster := k8scluster.NewOperator(user).WithQueryParma(parma)
		total, clusters, err := cluster.List(nil)
		if utils.HandleError(c, err) {
			return
		}

		// 真实集群kubeconfig脱敏
		for i := range clusters {
			clusters[i].KubeConfig = ""
		}
		data := map[string]interface{}{"list": clusters, "total": total}
		utils.HandleResult(c, &data)
	}
}

// GetK8sClustersByParams s2s接口根据ClusterName、ClusterId、PageSize、PageNum 选填参数获取集群信息列表
func (m *K8sCluster) GetK8sClustersByParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterName, clusterID := c.Query("clusterName"), utils.GetQueryIntItem("clusterId", "-1", c)
		pageSize, pageNum := utils.GetPageInfo(c)
		user := utils.PickUserInSession(c)
		authed := utils.PickUserAuthorized(c)
		parma := utils.NewQueryParma(pageSize, pageNum, true, authed)
		cluster := k8scluster.NewOperator(user).WithQueryParma(parma)
		// assemble attach params with table field => value
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
