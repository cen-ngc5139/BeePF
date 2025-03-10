package models

import (
	"time"
)

// TaskDB 任务数据库模型
type TaskDB struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"column:name" json:"name"`
	Description    string    `gorm:"column:description;type:text" json:"description"`
	ComponentID    uint64    `gorm:"column:component_id;index" json:"component_id"`
	ComponentName  string    `gorm:"column:component_name" json:"component_name"`
	Step           int       `gorm:"column:step;comment:任务步骤" json:"step"`
	Status         int       `gorm:"column:status;comment:任务状态" json:"status"`
	Error          string    `gorm:"column:error;type:text" json:"error"`
	Deleted        uint8     `gorm:"column:deleted;default:0" json:"deleted"`
	Creator        string    `gorm:"column:creator" json:"creator"`
	CreatedTime    time.Time `gorm:"column:created_time;autoCreateTime" json:"created_time"`
	LastUpdateTime time.Time `gorm:"column:last_update_time;autoUpdateTime" json:"last_update_time"`

	// 关联关系
	ProgStatuses []TaskProgStatusDB `gorm:"foreignKey:TaskID" json:"prog_statuses"`
}

// TableName 指定表名
func (TaskDB) TableName() string {
	return "beepf.task"
}

// ToTask 将数据库模型转换为业务模型
func (t *TaskDB) ToTask() *Task {
	task := &Task{
		ID:            t.ID,
		Name:          t.Name,
		Description:   t.Description,
		ComponentID:   t.ComponentID,
		ComponentName: t.ComponentName,
		Step:          TaskStep(t.Step),
		Status:        TaskStatus(t.Status),
		Error:         t.Error,
		CreatedAt:     t.CreatedTime,
		UpdatedAt:     t.LastUpdateTime,
	}

	// 转换 ProgStatuses
	if len(t.ProgStatuses) > 0 {
		task.ProgStatus = make([]ComProgStatus, len(t.ProgStatuses))
		for i, ps := range t.ProgStatuses {
			task.ProgStatus[i] = *ps.ToComProgStatus()
		}
	}

	return task
}

// TaskProgStatusDB 任务程序状态数据库模型
type TaskProgStatusDB struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID         uint64    `gorm:"column:task_id;index" json:"task_id"`
	ComponentID    uint64    `gorm:"column:component_id;index" json:"component_id"`
	ComponentName  string    `gorm:"column:component_name" json:"component_name"`
	ProgramID      uint64    `gorm:"column:program_id;index" json:"program_id"`
	ProgramName    string    `gorm:"column:program_name" json:"program_name"`
	Status         int       `gorm:"column:status;comment:状态" json:"status"`
	Error          string    `gorm:"column:error;type:text" json:"error"`
	Deleted        uint8     `gorm:"column:deleted;default:0" json:"deleted"`
	CreatedTime    time.Time `gorm:"column:created_time;autoCreateTime" json:"created_time"`
	LastUpdateTime time.Time `gorm:"column:last_update_time;autoUpdateTime" json:"last_update_time"`
}

// TableName 指定表名
func (TaskProgStatusDB) TableName() string {
	return "beepf.task_program_status"
}

// ToComProgStatus 将数据库模型转换为业务模型
func (c *TaskProgStatusDB) ToComProgStatus() *ComProgStatus {
	return &ComProgStatus{
		ID:            c.ID,
		TaskID:        c.TaskID,
		ComponentID:   c.ComponentID,
		ComponentName: c.ComponentName,
		ProgramID:     c.ProgramID,
		ProgramName:   c.ProgramName,
		Status:        TaskStatus(c.Status),
		Error:         c.Error,
		CreatedAt:     c.CreatedTime,
		UpdatedAt:     c.LastUpdateTime,
	}
}

// FromTask 从业务模型创建数据库模型
func TaskDBFromTask(task *Task) *TaskDB {
	taskDB := &TaskDB{
		ID:             task.ID,
		Name:           task.Name,
		Description:    task.Description,
		ComponentID:    task.ComponentID,
		ComponentName:  task.ComponentName,
		Step:           int(task.Step),
		Status:         int(task.Status),
		Error:          task.Error,
		CreatedTime:    task.CreatedAt,
		LastUpdateTime: task.UpdatedAt,
	}

	// 转换 ProgStatus
	if len(task.ProgStatus) > 0 {
		taskDB.ProgStatuses = make([]TaskProgStatusDB, len(task.ProgStatus))
		for i, ps := range task.ProgStatus {
			taskDB.ProgStatuses[i] = TaskProgStatusDB{
				ID:             ps.ID,
				TaskID:         ps.TaskID,
				ComponentID:    ps.ComponentID,
				ComponentName:  ps.ComponentName,
				ProgramID:      ps.ProgramID,
				ProgramName:    ps.ProgramName,
				Status:         int(ps.Status),
				Error:          ps.Error,
				CreatedTime:    ps.CreatedAt,
				LastUpdateTime: ps.UpdatedAt,
			}
		}
	}

	return taskDB
}
