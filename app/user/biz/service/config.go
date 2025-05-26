package service

import (
	"crypto/rand"
	"encoding/base64"
	user_service "tiktok_e-commerce/rpc_gen/kitex_gen/user"

	"golang.org/x/crypto/bcrypt"
)

// 把状态码简化
func buildResponse(message string, status bool) *user_service.ResponseStatus {
	return &user_service.ResponseStatus{
		Message: message,
		Status:  status,
	}
}

// 生成随机盐值
func generateSalt(length int) (string , error) {
	salt := make([]byte, length)
	_ , err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// 密码加盐并哈希
func hashPasswordWithSalt(password string, salt string) (string, error) {
	// 将盐值与密码结合
	saltedPassword := password + salt
	// 使用 bcrypt 进行哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
