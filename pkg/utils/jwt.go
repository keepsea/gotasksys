// pkg/utils/jwt.go

package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateToken 为指定用户生成JWT
func GenerateToken(userID uuid.UUID, userRole string, jwtSecret string, jwtExpirationHours int) (string, error) {
	// 设置Token的声明(Claims)
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    userRole,
		"exp":     time.Now().Add(time.Hour * time.Duration(jwtExpirationHours)).Unix(), // 过期时间
		"iat":     time.Now().Unix(),                                                    // 签发时间
	}

	// 创建一个 signing token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥进行签名，生成最终的token字符串
	signedToken, err := token.SignedString([]byte(jwtSecret))

	return signedToken, err
}
