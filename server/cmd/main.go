package main

import (
	"flag"
	"log"
	"time"

	"github.com/cen-ngc5139/BeePF/server/internal/database"
	"github.com/cen-ngc5139/BeePF/server/internal/metrics"
	"go.uber.org/zap"

	"github.com/cen-ngc5139/BeePF/server/conf"
	"github.com/cen-ngc5139/BeePF/server/router"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := flag.String("config", "", "configuration file")
	flag.Parse()
	conf.ParseConfig(*cfg, true)
	database.Setup()

	gin.SetMode(gin.DebugMode)

	stop := make(chan struct{})
	defer close(stop)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panicf("Start server failed, error :%v", err)
	}

	collector, err := metrics.NewNodeMetricsCollector(time.Second, logger)
	if err != nil {
		log.Panicf("Start server failed, error :%v", err)
	}
	defer collector.Stop()

	s := router.NewServer()
	s.SetMetricsCollector(collector)
	if err := s.Start(); err != nil {
		log.Panicf("Start server failed, error :%v", err)
	}
}
