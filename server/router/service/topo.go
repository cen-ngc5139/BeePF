package service

import (
	"errors"
	"net/http"

	"github.com/cen-ngc5139/BeePF/server/internal/operator/observability"
	"github.com/cen-ngc5139/BeePF/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Topo struct{}

func (t *Topo) Topo() gin.HandlerFunc {
	return func(c *gin.Context) {
		topoOp := observability.NewTopo()
		topo, err := topoOp.GetTopo()
		if utils.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, topo)
	}
}

func (t *Topo) Prog() gin.HandlerFunc {
	return func(c *gin.Context) {
		topoOp := observability.NewTopo()
		progs, err := topoOp.ListProgs()
		if utils.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, progs)
	}
}

func (t *Topo) ProgDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		progID := c.Param("progId")
		if progID == "" {
			utils.HandleError(c, errors.New("progId is required"))
			return
		}

		topoOp := observability.NewTopo()
		detail, err := topoOp.GetProgDetail(progID)
		if utils.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, detail)
	}
}

func (t *Topo) ProgDump() gin.HandlerFunc {
	return func(c *gin.Context) {
		progID := c.Param("progId")
		if progID == "" {
			utils.HandleError(c, errors.New("progId is required"))
			return
		}
		topoOp := observability.NewTopo()
		dump, err := topoOp.GetProgDump(progID)
		if utils.HandleError(c, err) {
			return
		}

		c.Data(http.StatusOK, "text/plain", dump)
	}
}
