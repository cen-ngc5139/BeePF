package service

import (
	"strconv"

	"github.com/cen-ngc5139/BeePF/server/internal/operator/component"
	"github.com/cen-ngc5139/BeePF/server/internal/operator/task"
	"github.com/cen-ngc5139/BeePF/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Task struct{}

// Create 创建任务
func (t *Task) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取组件ID
		componentId := c.Param("componentId")
		id, err := strconv.ParseUint(componentId, 10, 64)
		if utils.HandleError(c, err) {
			return
		}

		// 获取组件信息
		componentOp := component.NewOperator()
		comp, err := componentOp.Get(id)
		if utils.HandleError(c, err) {
			return
		}

		// 创建并运行任务
		taskOp := task.NewOperator()
		createdTask, err := taskOp.CreateAndRunTask(comp)
		if utils.HandleError(c, err) {
			return
		}

		utils.HandleResult(c, createdTask)
	}
}

// List 获取任务列表
func (t *Task) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageSize, pageNum := utils.GetPageInfo(c)
		parma := utils.NewQueryParma(pageSize, pageNum)

		taskOp := task.NewOperator().WithQueryParma(parma)
		total, tasks, err := taskOp.TaskStore.ListTasks(parma)
		if utils.HandleError(c, err) {
			return
		}

		data := map[string]interface{}{"list": tasks, "total": total}
		utils.HandleResult(c, &data)
	}
}

// Get 获取任务详情
func (t *Task) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取任务ID
		taskId := c.Param("taskId")
		id, err := strconv.ParseUint(taskId, 10, 64)
		if utils.HandleError(c, err) {
			return
		}

		// 获取任务信息
		taskOp := task.NewOperator()
		taskInfo, err := taskOp.TaskStore.GetTask(id)
		if utils.HandleError(c, err) {
			return
		}

		utils.HandleResult(c, taskInfo)
	}
}

// Stop 停止任务
func (t *Task) Stop() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取任务ID
		taskId := c.Param("taskId")
		id, err := strconv.ParseUint(taskId, 10, 64)
		if utils.HandleError(c, err) {
			return
		}

		// 停止任务
		taskOp := task.NewOperator()
		err = taskOp.StopTask(id)
		if utils.HandleError(c, err) {
			return
		}

		utils.HandleResult(c, nil)
	}
}

// Running 获取正在运行的任务
func (t *Task) Running() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskOp := task.NewOperator()
		runningTasks := taskOp.GetRunningTasks()

		data := map[string]interface{}{"list": runningTasks, "total": len(runningTasks)}
		utils.HandleResult(c, &data)
	}
}

// Metrics 获取任务指标
func (t *Task) Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskId := c.Param("taskId")
		id, err := strconv.ParseUint(taskId, 10, 64)
		if utils.HandleError(c, err) {
			return
		}

		taskOp := task.NewOperator()
		metrics, err := taskOp.GetTaskMetrics(id)
		if utils.HandleError(c, err) {
			return
		}

		utils.HandleResult(c, metrics)
	}
}
