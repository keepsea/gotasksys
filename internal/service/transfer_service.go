// internal/service/transfer_service.go
package service

import (
	"errors"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"

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

func RespondToTransferService(transferID, responderID uuid.UUID, action string) error {
	transfer, err := repository.FindTransferByID(transferID)
	if err != nil {
		return errors.New("transfer request not found")
	}
	if transfer.Status != "pending" {
		return errors.New("transfer request is no longer pending")
	}
	if transfer.ToUserID != responderID {
		return errors.New("permission denied: you are not the designated recipient of this transfer")
	}

	if action == "accept" {
		task, _ := repository.FindTaskByID(transfer.TaskID)
		newEffort := task.OriginalEffort - transfer.EffortSpentByInitiator
		if newEffort < 0 {
			newEffort = 0
		}

		updates := map[string]interface{}{
			"status":      "in_progress",
			"assignee_id": transfer.ToUserID,
			"effort":      newEffort,
		}
		if err := repository.UpdateTaskFields(transfer.TaskID, updates); err != nil {
			return err
		}
		return repository.UpdateTransferStatus(transferID, "accepted")
	} else if action == "reject" {
		updates := map[string]interface{}{"status": "in_progress"}
		if err := repository.UpdateTaskFields(transfer.TaskID, updates); err != nil {
			return err
		}
		return repository.UpdateTransferStatus(transferID, "rejected")
	}
	return errors.New("invalid action")
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
