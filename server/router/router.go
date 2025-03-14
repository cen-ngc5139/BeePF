package router

import (
	"github.com/cen-ngc5139/BeePF/server/router/service"
	"github.com/gin-gonic/gin"
)

var (
	clusterService   = &service.Cluster{}
	componentService = &service.Component{}
	taskService      = &service.Task{}
	topoService      = &service.Topo{}
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
		v1.GET("/component", componentService.List())
		v1.GET("/component/:componentId", componentService.Get())
		v1.POST("/component", componentService.Create())
		v1.POST("/component/upload", componentService.Upload())
		// v1.PUT("/component/:componentId", componentService.Update())
		v1.DELETE("/component/:componentId", componentService.Delete())

		// 任务管理相关接口
		v1.GET("/task", taskService.List())
		v1.GET("/task/:taskId", taskService.Get())
		v1.POST("/task/component/:componentId", taskService.Create())
		v1.GET("/task/running", taskService.Running())
		v1.POST("/task/:taskId/stop", taskService.Stop())
		v1.GET("/task/:taskId/metrics", taskService.Metrics())

		// 可观测相关接口
		v1.GET("/observability/topo", topoService.Topo())
	}

	return s.router
}
