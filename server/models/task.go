package models

import (
	"time"

	"github.com/pkg/errors"
)

type Task struct {
	ID            uint64          `json:"id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	ComponentID   uint64          `json:"component_id"`
	ComponentName string          `json:"component_name"`
	Step          TaskStep        `json:"step"`
	Status        TaskStatus      `json:"status"`
	Error         string          `json:"error"`
	ProgStatus    []ComProgStatus `json:"prog_status"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type TaskStep int

const (
	TaskStepInit TaskStep = iota
	TaskStepLoad
	TaskStepStart
	TaskStepStats
	TaskStepMetrics
	TaskStepStop
)

type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusRunning
	TaskStatusSuccess
	TaskStatusFailed
)

type ComProgStatus struct {
	ID            uint64     `json:"id"`
	TaskID        uint64     `json:"task_id"`
	ComponentID   uint64     `json:"component_id"`
	ComponentName string     `json:"component_name"`
	ProgramID     uint64     `json:"program_id"`
	ProgramName   string     `json:"program_name"`
	Status        TaskStatus `json:"status"`
	Error         string     `json:"error"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (t *Task) Validate() error {
	if t.ComponentID == 0 {
		return errors.New("component_id is required")
	}

	if t.ComponentName == "" {
		return errors.New("name is required")
	}

	return nil
}
