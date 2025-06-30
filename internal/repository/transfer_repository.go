// internal/repository/transfer_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"

	"github.com/google/uuid"
)

func CreateTransfer(transfer *model.TaskTransfer) error {
	return config.DB.Create(transfer).Error
}

func FindTransferByID(id uuid.UUID) (model.TaskTransfer, error) {
	var transfer model.TaskTransfer
	err := config.DB.First(&transfer, "id = ?", id).Error
	return transfer, err
}

func UpdateTransferStatus(id uuid.UUID, status string) error {
	return config.DB.Model(&model.TaskTransfer{}).Where("id = ?", id).Update("status", status).Error
}
