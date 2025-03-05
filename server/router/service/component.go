package service

import (
	"errors"

	"github.com/cen-ngc5139/BeePF/server/internal/operator/component"
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/cen-ngc5139/BeePF/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Component struct{}

func (ct *Component) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		basic := &models.Component{}
		if err := c.BindJSON(basic); utils.HandleError(c, err) {
			return
		}

		handler := component.NewOperator().WithComponent(basic)
		if err := handler.Create(); utils.HandleError(c, err) {
			return
		}

		data := &map[string]interface{}{}
		utils.HandleResult(c, data)
	}
}

func (ct *Component) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := utils.GetParamIntItem("componentId", c)

		if id == 0 {
			utils.ResponseErr(c, errors.New("组件编号不合法"))
			return
		}

		operator := component.NewOperator()
		component, err := operator.Get(uint64(id))
		if utils.HandleError(c, err) {
			return
		}

		data := &map[string]interface{}{"component": component}
		utils.HandleResult(c, data)
	}
}

func (ct *Component) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageSize, pageNum := utils.GetPageInfo(c)
		parma := utils.NewQueryParma(pageSize, pageNum)

		component := component.NewOperator().WithQueryParma(parma)
		total, components, err := component.List()
		if utils.HandleError(c, err) {
			return
		}

		data := map[string]interface{}{"list": components, "total": total}
		utils.HandleResult(c, &data)
	}
}
