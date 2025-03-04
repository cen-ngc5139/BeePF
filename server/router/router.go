package router

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) initRouter() *gin.Engine {
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/ping", Ping)

	}

	return s.router
}
