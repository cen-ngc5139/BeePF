package router

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"

	"github.com/cen-ngc5139/BeePF/server/conf"
	"github.com/cen-ngc5139/BeePF/server/internal/metrics"
	"github.com/cen-ngc5139/BeePF/server/pkg/utils"
	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
)

var (
	InitCompleted bool
)

type Server struct {
	router  *gin.Engine
	server  http.Server
	metrics *metrics.NodeMetricsCollector
}

func Ping(c *gin.Context) {
	utils.ResponseOk(c, gin.H{
		"message": "pong",
	})
}

func NewServer(middleware ...gin.HandlerFunc) *Server {
	r := gin.Default()
	r.GET("/ping", Ping)

	InitProbe(r)
	InitPrometheusMetrics(r)
	pprof.Register(r, "pprof")

	r.Use(middleware...)
	return &Server{
		router: r,
		server: http.Server{
			Addr:    conf.Config().Http.Listen,
			Handler: r,
		},
	}
}

func (s *Server) SetMetricsCollector(collector *metrics.NodeMetricsCollector) {
	s.metrics = collector
}

func (s *Server) Start() error {
	// register api
	s.initRouter()

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("Listen: %s\n", err)
		}
	}()

	InitCompleted = true

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 2)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	return nil
}

func InitProbe(r *gin.Engine) {
	r.GET("/readiness", readiness())
	r.GET("/liveness", liveness())
}

func liveness() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.HandleResult(c, nil)
	}
}

func readiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		if InitCompleted {
			utils.HandleResult(c, nil)
			return
		}

		utils.ResponseProbeError(c, errors.New("server not ready"))
	}
}
