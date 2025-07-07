// internal/service/transfer_service.go
package service

import (
	"errors"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"gotasksys/pkg/apierror"
	"log"

	"github.com/google/uuid"
)

func InitiateTransferService(taskID uint, initiatorID, newAssigneeID uuid.UUID, effortSpent int) (model.TaskTransfer, error) {
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return model.TaskTransfer{}, errors.New("task not found")
	}
	if task.Status != "in_progress" {
		return model.TaskTransfer{}, errors.New("task is not in progress")
	}
	if task.AssigneeID == nil || *task.AssigneeID != initiatorID {
		return model.TaskTransfer{}, errors.New("permission denied: you are not the current assignee")
	}

	transfer := model.TaskTransfer{
		TaskID:                 taskID,
		FromUserID:             initiatorID,
		ToUserID:               newAssigneeID,
		EffortSpentByInitiator: effortSpent,
	}

	if err := repository.CreateTransfer(&transfer); err != nil {
		return model.TaskTransfer{}, err
	}

	updates := map[string]interface{}{"status": "pending_transfer"}
	if err := repository.UpdateTaskFields(taskID, updates); err != nil {
		return model.TaskTransfer{}, err
	}

	return transfer, nil
}

// RespondToTransferService 封装了响应转交请求（接受/拒绝）的业务逻辑 (最终级联版)
func RespondToTransferService(transferID, respondentID uuid.UUID, action string) error {
	// 1. 查找转交记录并校验
	transfer, err := repository.FindTransferByID(transferID)
	if err != nil {
		return apierror.ErrTransferNotFound
	}
	if transfer.Status != "pending" {
		return apierror.ErrTransferStatusConflict
	}
	// 权限校验：只有被指定的接收人才能响应
	if transfer.ToUserID != respondentID {
		return apierror.ErrPermissionDenied
	}

	// 2. 根据action字符串来处理
	if action == "accept" {
		// --- 【接受转交】的完整逻辑 ---

		// a. 将转交记录的状态更新为 'accepted'
		if err := repository.UpdateTransferStatus(transferID, "accepted"); err != nil {
			return err
		}

		// b. 更新主任务的负责人、状态和重新计算的工时
		task, err := repository.FindTaskByID(transfer.TaskID)
		if err != nil {
			return err
		}
		newEffort := task.OriginalEffort - transfer.EffortSpentByInitiator
		if newEffort < 0 {
			newEffort = 0
		}
		mainTaskUpdates := map[string]interface{}{
			"assignee_id": respondentID,
			"status":      "in_progress",
			"effort":      newEffort,
		}
		if err := repository.UpdateTaskFields(transfer.TaskID, mainTaskUpdates); err != nil {
			return err
		}

		// c. 【核心修正】级联转交子任务
		// 将所有隶属于该主任务、且负责人是原负责人(FromUserID)的子任务，一并转交给新负责人(respondentID)
		err = repository.BatchUpdateSubtasksAssignee(transfer.TaskID, transfer.FromUserID, respondentID)
		if err != nil {
			// 即使这一步失败，主任务的转交也已完成，所以我们只记录日志而不返回阻塞性错误
			log.Printf("Warning: Failed to cascade subtask assignee update for parent task %d. Error: %v", transfer.TaskID, err)
		}

	} else if action == "reject" {
		// --- 【拒绝转交】的逻辑 ---
		// 将转交记录状态更新为 'rejected'
		if err := repository.UpdateTransferStatus(transferID, "rejected"); err != nil {
			return err
		}
		// 将原任务的状态恢复为 'in_progress'
		updates := map[string]interface{}{"status": "in_progress"}
		if err := repository.UpdateTaskFields(transfer.TaskID, updates); err != nil {
			log.Printf("Warning: Transfer status set to 'rejected' but failed to update task %d status to 'in_progress': %v", transfer.TaskID, err)
			return err
		}
		return nil

	} else {
		return errors.New("invalid action specified")
	}

	return nil
}

// CancelTransferService 封装了取消转交的业务逻辑
func CancelTransferService(transferID, initiatorID uuid.UUID) error {
	// 1. 查找转交记录并校验
	transfer, err := repository.FindTransferByID(transferID)
	if err != nil {
		return errors.New("transfer request not found")
	}
	if transfer.Status != "pending" {
		return errors.New("transfer request is no longer pending and cannot be cancelled")
	}
	// 权限校验：只有发起人自己才能取消
	if transfer.FromUserID != initiatorID {
		return errors.New("permission denied: you are not the initiator of this transfer")
	}

	// 2. 将转交记录的状态更新为 'cancelled'
	if err := repository.UpdateTransferStatus(transferID, "cancelled"); err != nil {
		return err
	}

	// 3. 将原任务的状态恢复为 'in_progress'
	updates := map[string]interface{}{"status": "in_progress"}
	return repository.UpdateTaskFields(transfer.TaskID, updates)
}
