package main

import (
	"flag"
	"github.com/cen-ngc5139/BeePF/server/internal/database"
	"log"

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

	s := router.NewServer()
	if err := s.Start(); err != nil {
		log.Panicf("Start server failed, error :%v", err)
	}
}
