package task

import (
	"time"

	"github.com/cen-ngc5139/BeePF/server/internal/database"
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/cen-ngc5139/BeePF/server/pkg/utils"
	"gorm.io/gorm"
)

type Store struct {
}

func (s *Store) CreateTask(task *models.Task) (*models.Task, error) {
	var taskDB *models.TaskDB
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		taskDB = &models.TaskDB{
			ID:            task.ID,
			Name:          task.Name,
			Status:        int(task.Status),
			Step:          int(task.Step),
			Description:   task.Description,
			ComponentID:   task.ComponentID,
			ComponentName: task.ComponentName,
			Error:         task.Error,
			CreatedTime:   task.CreatedAt,
		}

		if err := tx.Create(taskDB).Error; err != nil {
			return err
		}

		// 更新任务ID
		task.ID = taskDB.ID

		for _, program := range task.ProgStatus {
			programDB := &models.TaskProgStatusDB{
				TaskID:        taskDB.ID,
				ComponentID:   program.ComponentID,
				ComponentName: program.ComponentName,
				ProgramID:     program.ProgramID,
				ProgramName:   program.ProgramName,
				Status:        int(program.Status),
				Error:         program.Error,
				CreatedTime:   program.CreatedAt,
			}

			if err := tx.Create(programDB).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 单独查询任务及其关联数据
	result := database.DB.Preload("ProgStatuses").First(taskDB, taskDB.ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return taskDB.ToTask(), nil
}

func (s *Store) GetTask(id uint64) (*models.Task, error) {
	var taskDB models.TaskDB
	result := database.DB.Preload("ProgStatuses").First(&taskDB, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return taskDB.ToTask(), nil
}

func (s *Store) ListTasks(query *utils.Query) (total int64, tasks []*models.Task, err error) {
	var tasksDB []models.TaskDB

	// 计算总数
	result := database.DB.Model(&models.TaskDB{}).Count(&total)
	if result.Error != nil {
		return 0, nil, result.Error
	}

	db := database.DB.Preload("ProgStatuses")

	if query != nil && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize
		if offset < 0 {
			offset = 0
		}
		db = db.Limit(query.PageSize).Offset(offset)
	}

	// 按ID倒序排序，确保最新的任务显示在前面
	result = db.Order("id DESC").Find(&tasksDB)
	if result.Error != nil {
		return 0, nil, result.Error
	}

	// 转换为业务模型
	tasks = make([]*models.Task, len(tasksDB))
	for i, taskDB := range tasksDB {
		tasks[i] = taskDB.ToTask()
	}

	return total, tasks, nil
}

func (s *Store) UpdateTask(task *models.Task) error {
	// 首先获取现有的任务记录
	var existingTask models.TaskDB
	if err := database.DB.First(&existingTask, task.ID).Error; err != nil {
		return err
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 只更新需要更新的字段，保留其他字段不变
		taskDB := existingTask

		// 更新需要修改的字段
		taskDB.Status = int(task.Status)
		taskDB.Step = int(task.Step)
		taskDB.Error = task.Error

		// 确保更新时间有效
		if !task.UpdatedAt.IsZero() {
			taskDB.LastUpdateTime = task.UpdatedAt
		} else {
			taskDB.LastUpdateTime = time.Now()
		}

		// 使用Select指定要更新的字段，避免更新零值
		if err := tx.Model(&taskDB).Select("status", "step", "error", "last_update_time").Updates(taskDB).Error; err != nil {
			return err
		}

		// 更新程序状态
		if len(task.ProgStatus) > 0 {
			for _, program := range task.ProgStatus {
				// 首先检查程序状态记录是否存在
				var existingProgram models.TaskProgStatusDB
				err := tx.Where("task_id = ? AND program_id = ?", task.ID, program.ProgramID).First(&existingProgram).Error

				if err != nil {
					return err
				}

				// 如果记录存在，更新记录
				existingProgram.Status = int(program.Status)
				existingProgram.Error = program.Error

				// 确保更新时间有效
				if !program.UpdatedAt.IsZero() {
					existingProgram.LastUpdateTime = program.UpdatedAt
				} else {
					existingProgram.LastUpdateTime = time.Now()
				}

				if err := tx.Model(&existingProgram).Select("status", "error", "last_update_time").Updates(existingProgram).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *Store) DeleteTask(task *models.Task) error {
	return nil
}
