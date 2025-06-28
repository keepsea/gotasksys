// pkg/utils/password.go

package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14是加密成本，数值越高越安全但越慢
	return string(bytes), err
}

// CheckPasswordHash 检查密码和哈希值是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
