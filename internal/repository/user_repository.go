// internal/repository/user_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"

	"github.com/google/uuid"
)

func FindUserByUsername(username string) (model.User, error) {
	var user model.User
	result := config.DB.Where("username = ?", username).First(&user)
	return user, result.Error
}

func FindUserByID(id uuid.UUID) (model.User, error) {
	var user model.User
	result := config.DB.First(&user, "id = ?", id)
	return user, result.Error
}

func CreateUser(user *model.User) error {
	result := config.DB.Create(user)
	return result.Error
}
