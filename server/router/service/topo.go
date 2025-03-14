package service

import (
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
