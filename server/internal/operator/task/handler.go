package task

import (
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/pkg/errors"
)

func (o *Operator) GetTask(id uint64) (task *models.Task, err error) {
	task, err = o.TaskStore.GetTask(id)
	if err != nil {
		err = errors.Wrap(err, "获取任务失败")
		return
	}

	return
}

func (o *Operator) ListTask() (total int64, tasks []*models.Task, err error) {
	total, tasks, err = o.TaskStore.ListTasks(o.QueryParma)
	if err != nil {
		err = errors.Wrap(err, "获取任务列表失败")
		return
	}

	return
}

func (o *Operator) UpdateTask(task *models.Task) (err error) {
	err = o.TaskStore.UpdateTask(task)
	if err != nil {
		err = errors.Wrap(err, "更新任务失败")
		return
	}

	return
}

func (o *Operator) DeleteTask(id uint64) (err error) {
	err = o.TaskStore.DeleteTask(o.Task)
	if err != nil {
		err = errors.Wrap(err, "删除任务失败")
		return
	}

	return
}
