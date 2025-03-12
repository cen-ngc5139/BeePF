package cache

import "sync"

// 用于存储正在运行的任务
// key: task id(uint64)
// value: task status(models.RunningTask)
var TaskRunningStore *sync.Map

func init() {
	TaskRunningStore = &sync.Map{}
}
