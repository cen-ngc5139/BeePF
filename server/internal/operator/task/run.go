package task

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	loader "github.com/cen-ngc5139/BeePF/loader/lib/src/cli"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"github.com/cen-ngc5139/BeePF/server/internal/cache"
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// CreateAndRunTask 创建并运行任务
func (o *Operator) CreateAndRunTask(component *models.Component) (*models.Task, error) {
	// 创建任务记录
	task := &models.Task{
		ComponentID:   uint64(component.Id),
		ComponentName: component.Name,
		Name:          component.Name + "-" + uuid.New().String()[:8],
		Description:   "运行组件 " + component.Name,
		Status:        models.TaskStatusPending,
		Step:          models.TaskStepInit,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 为每个程序创建状态记录
	if component.Programs != nil {
		task.ProgStatus = make([]models.ComProgStatus, len(component.Programs))
		for i, prog := range component.Programs {
			task.ProgStatus[i] = models.ComProgStatus{
				ComponentID:   uint64(component.Id),
				ComponentName: component.Name,
				ProgramID:     uint64(prog.Id),
				ProgramName:   prog.Name,
				Status:        models.TaskStatusPending,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
		}
	}

	// 保存任务到数据库
	createdTask, err := o.TaskStore.CreateTask(task)
	if err != nil {
		return nil, errors.Wrap(err, "创建任务失败")
	}

	// 在协程中运行任务
	go o.RunComponentAsync(createdTask, component)

	return createdTask, nil
}

// RunComponentAsync 异步运行组件
func (o *Operator) RunComponentAsync(task *models.Task, component *models.Component) {
	// 创建上下文，用于取消任务
	ctx, cancel := context.WithCancel(context.Background())

	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		task.Status = models.TaskStatusFailed
		task.Error = "初始化日志失败: " + err.Error()
		o.TaskStore.UpdateTask(task)
		cancel()
		return
	}
	defer logger.Sync()

	// 更新任务状态为运行中
	task.Status = models.TaskStatusRunning
	task.UpdatedAt = time.Now()
	err = o.TaskStore.UpdateTask(task)
	if err != nil {
		logger.Error("更新任务状态失败", zap.Error(err))
		cancel()
		return
	}

	// 创建并存储运行中的任务
	runningTask := &models.RunningTask{
		Task:       task,
		CancelFunc: cancel,
		Logger:     logger,
	}
	cache.TaskRunningStore.Store(task.ID, runningTask)

	// 确保在函数结束时清理资源
	defer func() {
		cache.TaskRunningStore.Delete(task.ID)
	}()

	// 配置BPF加载器
	config := &loader.Config{
		ObjectPath:  component.BinaryPath, // 这里应该使用组件的实际二进制路径
		Logger:      logger,
		PollTimeout: 100 * time.Millisecond,
		Properties: meta.Properties{
			Stats: &meta.Stats{
				Interval: 1 * time.Second,
				Handler:  metrics.NewDefaultHandler(logger),
			},
		},
	}

	bpfLoader := loader.NewBPFLoader(config)
	runningTask.BPFLoader = bpfLoader

	// 初始化BPF加载器
	task.Step = models.TaskStepInit
	task.UpdatedAt = time.Now()
	err = o.TaskStore.UpdateTask(task)
	if err != nil {
		logger.Error("更新任务步骤失败", zap.Error(err))
		return
	}

	err = bpfLoader.Init()
	if err != nil {
		logger.Error("初始化 BPF 加载器失败", zap.Error(err))
		task.Status = models.TaskStatusFailed
		task.Error = "初始化 BPF 加载器失败: " + err.Error()
		task.UpdatedAt = time.Now()
		if updateErr := o.TaskStore.UpdateTask(task); updateErr != nil {
			logger.Error("更新任务状态失败", zap.Error(updateErr))
		}
		return
	}

	cache.TaskRunningStore.Store(task.ID, runningTask)

	// 加载BPF程序
	task.Step = models.TaskStepLoad
	task.UpdatedAt = time.Now()
	err = o.TaskStore.UpdateTask(task)
	if err != nil {
		logger.Error("更新任务步骤失败", zap.Error(err))
		return
	}

	err = bpfLoader.Load()
	if err != nil {
		logger.Error("加载 BPF 程序失败", zap.Error(err))
		task.Status = models.TaskStatusFailed
		task.Error = "加载 BPF 程序失败: " + err.Error()
		task.UpdatedAt = time.Now()
		if updateErr := o.TaskStore.UpdateTask(task); updateErr != nil {
			logger.Error("更新任务状态失败", zap.Error(updateErr))
		}
		return
	}

	// 启动BPF程序
	task.Step = models.TaskStepStart
	task.UpdatedAt = time.Now()

	runningTask.BPFLoader = bpfLoader
	cache.TaskRunningStore.Store(task.ID, runningTask)

	// 更新程序状态
	progStatuses := make([]models.ComProgStatus, 0, len(task.ProgStatus))
	for _, prog := range task.ProgStatus {
		status, ok := bpfLoader.ProgAttachStatus[prog.ProgramName]
		if !ok {
			prog.Status = models.TaskStatusFailed
			prog.Error = "程序未找到"
		}

		prog.Status = models.TaskStatus(status.Status)
		prog.AttachID = status.AttachID
		prog.Error = status.Error
		progStatuses = append(progStatuses, prog)
	}

	task.ProgStatus = progStatuses

	err = o.TaskStore.UpdateTask(task)
	if err != nil {
		logger.Error("更新任务步骤失败", zap.Error(err))
		return
	}

	cache.TaskRunningStore.Store(task.ID, runningTask)

	if err := bpfLoader.Start(); err != nil {
		logger.Error("启动失败", zap.Error(err))
		task.Status = models.TaskStatusFailed
		task.Error = "启动失败: " + err.Error()
		task.UpdatedAt = time.Now()
		if updateErr := o.TaskStore.UpdateTask(task); updateErr != nil {
			logger.Error("更新任务状态失败", zap.Error(updateErr))
		}
		return
	}

	// 启动统计收集器
	task.Step = models.TaskStepStats
	task.UpdatedAt = time.Now()
	err = o.TaskStore.UpdateTask(task)
	if err != nil {
		logger.Error("更新任务步骤失败", zap.Error(err))
		return
	}

	cache.TaskRunningStore.Store(task.ID, runningTask)

	if err := bpfLoader.Stats(); err != nil {
		logger.Error("启动统计收集器失败", zap.Error(err))
		task.Status = models.TaskStatusFailed
		task.Error = "启动统计收集器失败: " + err.Error()
		task.UpdatedAt = time.Now()
		if updateErr := o.TaskStore.UpdateTask(task); updateErr != nil {
			logger.Error("更新任务状态失败", zap.Error(updateErr))
		}
		return
	}

	// 启动指标
	task.Step = models.TaskStepMetrics
	task.UpdatedAt = time.Now()
	err = o.TaskStore.UpdateTask(task)
	if err != nil {
		logger.Error("更新任务步骤失败", zap.Error(err))
		return
	}

	if err := bpfLoader.Metrics(); err != nil {
		logger.Error("启动指标失败", zap.Error(err))
		task.Status = models.TaskStatusFailed
		task.Error = "启动指标失败: " + err.Error()
		task.UpdatedAt = time.Now()
		if updateErr := o.TaskStore.UpdateTask(task); updateErr != nil {
			logger.Error("更新任务状态失败", zap.Error(updateErr))
		}
		return
	}

	cache.TaskRunningStore.Store(task.ID, runningTask)

	// 等待取消信号或上下文取消
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logger.Info("任务被取消")
	case <-sigChan:
		logger.Info("收到系统信号，正常关闭")
	}

	bpfLoader.Stop()

	// 停止BPF程序
	task.Step = models.TaskStepStop
	task.Status = models.TaskStatusSuccess
	task.UpdatedAt = time.Now()
	if err := o.TaskStore.UpdateTask(task); err != nil {
		logger.Error("更新任务状态失败", zap.Error(err))
	}

	logger.Info("任务完成")
}

// StopTask 停止正在运行的任务
func (o *Operator) StopTask(taskID uint64) error {
	runningTask, exists := cache.TaskRunningStore.Load(taskID)
	if !exists {
		return errors.New("任务不存在或已停止")
	}
	runningTask.(*models.RunningTask).CancelFunc()
	return nil
}

// GetRunningTasks 获取所有正在运行的任务
func (o *Operator) GetRunningTasks() []*models.Task {
	runningTasks := make([]*models.Task, 0)
	cache.TaskRunningStore.Range(func(key, value any) bool {
		runningTasks = append(runningTasks, value.(*models.RunningTask).Task)
		return true
	})
	return runningTasks
}
