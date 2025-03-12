package router

import (
	"github.com/cen-ngc5139/BeePF/server/internal/cache"
	"github.com/cen-ngc5139/BeePF/server/internal/metrics"
	"github.com/gin-gonic/gin"
)

func InitPrometheusMetrics(r *gin.Engine) {
	taskMetrics := metrics.NewTaskMetrics(cache.TaskRunningStore)
	r.GET("/metrics", taskMetrics.Handler())
}
