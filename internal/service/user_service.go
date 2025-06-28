// internal/service/user_service.go
package service

import (
	"errors"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"gotasksys/pkg/utils"
	"strings"

	"github.com/google/uuid"
)

func RegisterUser(input model.User) (uuid.UUID, error) {
	if len(input.PasswordHash) < 12 { // 这里的检查应该检查原始密码，我们稍后改进
		return uuid.Nil, errors.New("password must be at least 12 characters long")
	}

	hashedPassword, err := utils.HashPassword(input.PasswordHash) // 我们把原始密码临时存在PasswordHash字段
	if err != nil {
		return uuid.Nil, errors.New("failed to hash password")
	}
	input.PasswordHash = hashedPassword

	err = repository.CreateUser(&input)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return uuid.Nil, errors.New("username already exists")
		}
		return uuid.Nil, errors.New("failed to create user")
	}

	return input.ID, nil
}

func LoginUser(username, password string, jwtSecret string, jwtExpHours int) (string, error) {
	user, err := repository.FindUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, jwtSecret, jwtExpHours)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
