// internal/model/task.go (最终关联修正版)

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
	Status         string         `gorm:"type:varchar(50);not null" json:"status"`
	Priority       string         `gorm:"type:varchar(50);not null" json:"priority"`
	Effort         int            `json:"effort"`
	OriginalEffort int            `json:"original_effort"`
	Evaluation     datatypes.JSON `json:"evaluation,omitempty"`

	// --- 关联ID字段 ---
	CreatorID    uuid.UUID  `json:"creator_id"`
	TaskTypeID   uuid.UUID  `json:"task_type_id"` // 我们暂时不预加载这个，保持简单
	ReviewerID   *uuid.UUID `json:"reviewer_id,omitempty"`
	AssigneeID   *uuid.UUID `json:"assignee_id,omitempty"`
	ParentTaskID *uint      `gorm:"index" json:"parent_task_id,omitempty"`

	// --- 【V1.2 关键修正】定义GORM关联关系 ---
	// gorm:"foreignKey:..." 指定了当前模型中哪个字段是外键
	// gorm:"references:..." 指定了外键所引用的目标模型中的哪个字段（通常是主键'ID'）
	Creator  User  `gorm:"foreignKey:CreatorID;references:ID" json:"creator"`
	Reviewer *User `gorm:"foreignKey:ReviewerID;references:ID" json:"reviewer,omitempty"`
	Assignee *User `gorm:"foreignKey:AssigneeID;references:ID" json:"assignee,omitempty"`
	// ------------------------------------

	// 时间戳
	CreatedAt   time.Time  `json:"created_at"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	ClaimedAt   *time.Time `json:"claimed_at,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
