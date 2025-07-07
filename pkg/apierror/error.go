package apierror

// pkg/apierror/error.go

import "fmt"

// APIError 定义了我们自定义的、包含错误码的错误结构
type APIError struct {
	Code    int    // 错误码
	Message string // 错误信息 (供开发者调试)
}

// Error 方法让 APIError 结构体实现了标准的 error 接口
func (e *APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Message)
}

// NewAPIError 创建一个新的APIError实例
func NewAPIError(code int, message string) *APIError {
	return &APIError{Code: code, Message: message}
}

// --- 在这里统一定义我们系统中所有的业务逻辑错误 ---
var (
	// 用户与认证 (1xxx)
	ErrUserNotFound         = NewAPIError(1001, "user not found")
	ErrInvalidCredentials   = NewAPIError(1002, "invalid credentials")
	ErrUsernameExists       = NewAPIError(1003, "username already exists")
	ErrInvalidRole          = NewAPIError(1004, "invalid role specified")
	ErrPasswordTooShort     = NewAPIError(1005, "password must be at least 12 characters long")
	ErrOldPasswordIncorrect = NewAPIError(1006, "old password is incorrect")

	// 通用权限 (2xxx)
	ErrPermissionDenied = NewAPIError(2001, "permission denied")

	// 任务相关 (3xxx)
	ErrTaskNotFound          = NewAPIError(3001, "task not found")
	ErrTaskStatusConflict    = NewAPIError(3002, "task status conflict") // 泛指任务状态不正确
	ErrSubtaskEffortExceeds  = NewAPIError(3003, "total effort of subtasks cannot exceed parent task's original effort")
	ErrSubtaskDueDateExceeds = NewAPIError(3004, "subtask due date cannot be after the parent task's due date")
	ErrCompleteWithSubtasks  = NewAPIError(3005, "cannot complete main task: there are still incomplete subtasks")

	// 转交相关 (4xxx)
	ErrTransferNotFound       = NewAPIError(4001, "transfer request not found")
	ErrTransferStatusConflict = NewAPIError(4002, "transfer request is no longer pending")
)
