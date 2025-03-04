package router

import (
	"github.com/cen-ngc5139/BeePF/server/router/service"
	"github.com/gin-gonic/gin"
)

var (
	clusterService = &service.Cluster{}
)

func (s *Server) initRouter() *gin.Engine {
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/ping", Ping)

		// 集群管理相关接口
		// k8s_cluster 真实集群
		v1.GET("/cluster", clusterService.List())
		v1.GET("/cluster/:clusterId", clusterService.Get())
		v1.POST("/cluster", clusterService.Create())
		v1.PUT("/cluster/:clusterId", clusterService.Update())
		v1.DELETE("/cluster/:clusterId", clusterService.Delete())
		v1.GET("/clusterList", clusterService.GetClustersByParams())
		// 组件管理相关接口
		// components := v1.Group("/components")
		// {
		// // 获取组件列表
		// components.GET("", ListComponents)
		// // 获取单个组件详情
		// components.GET("/:id", GetComponent)
		// // 创建组件
		// components.POST("", CreateComponent)
		// // 更新组件
		// components.PUT("/:id", UpdateComponent)
		// // 删除组件
		// components.DELETE("/:id", DeleteComponent)
		// // 部署组件
		// components.POST("/:id/deploy", DeployComponent)
		// // 停止组件
		// components.POST("/:id/stop", StopComponent)
		// // 获取组件日志
		// components.GET("/:id/logs", GetComponentLogs)
		// // 获取组件指标
		// components.GET("/:id/metrics", GetComponentMetrics)
		// }
	}

	return s.router
}
