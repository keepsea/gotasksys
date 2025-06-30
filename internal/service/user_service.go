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

// LoginUser 验证用户凭据并返回用户对象和Token
// 注意：返回值从 (string, error) 变为了 (model.User, string, error)
func LoginUser(username, password string, jwtSecret string, jwtExpHours int) (model.User, string, error) {
	// 1. 查找用户
	user, err := repository.FindUserByUsername(username)
	if err != nil {
		// 如果出错，返回一个空的User对象、空token和错误
		return model.User{}, "", errors.New("user not found")
	}

	// 2. 验证密码
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return model.User{}, "", errors.New("invalid credentials")
	}

	// 3. 生成JWT
	token, err := utils.GenerateToken(user.ID, user.Role, jwtSecret, jwtExpHours)
	if err != nil {
		return model.User{}, "", errors.New("failed to generate token")
	}

	// 4. 成功时，返回完整的user对象、token字符串和nil错误
	return user, token, nil
}
