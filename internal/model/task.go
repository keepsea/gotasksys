// internal/model/task.go

package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Task 定义了任务的数据结构
type Task struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Title          string         `gorm:"type:varchar(255);not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	TaskTypeID     uuid.UUID      `json:"task_type_id"`
	Status         string         `gorm:"type:varchar(50);not null" json:"status"`
	Priority       string         `gorm:"type:varchar(50);not null" json:"priority"`
	Effort         int            `json:"effort"`
	OriginalEffort int            `json:"original_effort"`
	CreatorID      uuid.UUID      `json:"creator_id"`
	ReviewerID     *uuid.UUID     `json:"reviewer_id,omitempty"` // 使用指针以允许为NULL
	AssigneeID     *uuid.UUID     `json:"assignee_id,omitempty"` // 使用指针以允许为NULL
	ParentTaskID   *uint          `gorm:"index" json:"parent_task_id,omitempty"`
	Evaluation     datatypes.JSON `json:"evaluation,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	ApprovedAt     *time.Time     `json:"approved_at,omitempty"`
	ClaimedAt      *time.Time     `json:"claimed_at,omitempty"`
	DueDate        *time.Time     `json:"due_date,omitempty"`
	CompletedAt    *time.Time     `json:"completed_at,omitempty"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
