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
