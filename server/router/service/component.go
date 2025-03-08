package service

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"

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

func (ct *Component) Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取上传的文件
		fileHeader, err := c.FormFile("binary")
		if err != nil {
			utils.ResponseErr(c, err)
			return
		}

		// 从表单中获取 JSON 数据
		jsonData := c.PostForm("data")
		if jsonData == "" {
			utils.ResponseErr(c, errors.New("missing component data"))
			return
		}

		// 解析 JSON 数据
		basic := &models.Component{}
		if err := json.Unmarshal([]byte(jsonData), basic); err != nil {
			utils.ResponseErr(c, err)
			return
		}

		// 打开文件
		file, err := fileHeader.Open()
		if err != nil {
			utils.ResponseErr(c, err)
			return
		}
		defer file.Close()

		// 读取文件内容到字节数组
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			utils.ResponseErr(c, err)
			return
		}

		// 创建操作对象并设置二进制数据
		operator := component.NewOperator().WithComponent(basic)
		operator.Binary = fileBytes

		// 处理上传的二进制文件
		if err := operator.UploadBinary(); utils.HandleError(c, err) {
			return
		}

		utils.HandleResult(c, nil)
	}
}

// Delete 删除组件
func (ct *Component) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取组件ID
		componentId := c.Param("componentId")
		id, err := strconv.ParseUint(componentId, 10, 64)
		if utils.HandleError(c, err) {
			return
		}

		// 获取组件信息
		operator := component.NewOperator()
		componentObj, err := operator.Get(id)
		if utils.HandleError(c, err) {
			return
		}

		// 删除组件及其关联的程序和映射
		err = operator.WithComponent(componentObj).Delete()
		if utils.HandleError(c, err) {
			return
		}

		utils.HandleResult(c, nil)
	}
}
