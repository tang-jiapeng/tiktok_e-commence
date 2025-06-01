package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"tiktok_e-commerce/auth/biz/dal/redis"
	"tiktok_e-commerce/auth/conf"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ctx                        = context.Background()
	jwtSecret                  = []byte(conf.GetConf().Jwt.Secret)
	accessTokenExpireDuration  = time.Hour * 2
	refreshTokenExpireDuration = time.Hour * 24 * 7
)

const (
	// TokenValid 令牌有效
	TokenValid = iota
	// TokenInvalid 令牌不合法
	TokenInvalid
	// TokenExpired 令牌过期
	TokenExpired
)

func GenerateRefreshToken(userId int32) (string, error) {
	s, err := generateJWT(userId, refreshTokenExpireDuration)
	if err == nil {
		err = redis.RedisClient.Set(ctx, GetRefreshTokenKey(userId), s, refreshTokenExpireDuration).Err()
		if err != nil {
			return "", err
		}
	}
	return s, nil
}

func GenerateAccessToken(userId int32) (string, error) {
	s, err := generateJWT(userId, accessTokenExpireDuration)
	if err == nil {
		err = redis.RedisClient.Set(ctx, GetAccessTokenKey(userId), s, accessTokenExpireDuration).Err()
		if err != nil {
			return "", err
		}
	}
	return s, err
}

func generateJWT(userId int32, exp time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenStr string) (jwt.MapClaims, int) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	switch {
	case token.Valid:
		return token.Claims.(jwt.MapClaims), TokenValid
	case errors.Is(err, jwt.ErrTokenExpired), errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, TokenExpired
	case errors.Is(err, jwt.ErrTokenMalformed), errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, TokenInvalid
	default:
		return nil, TokenInvalid
	}
}

func saveRefreshToken(userId int32, refreshToken string) error {
	return redis.RedisClient.Set(ctx, GetRefreshTokenKey(userId), refreshToken, refreshTokenExpireDuration).Err()
}

func refreshAccessToken(refreshToken string) (string, bool) {
	// 解析refreshToken
	claims, status := ParseJWT(refreshToken)
	if status != TokenValid {
		return "", false
	}
	userId := claims["userId"].(int32)
	newAccessToken, err := generateJWT(userId, accessTokenExpireDuration)
	if err != nil {
		return "", false
	}
	err = redis.RedisClient.Set(ctx, GetAccessTokenKey(userId), newAccessToken, accessTokenExpireDuration).Err()
	if err != nil {
		return "", false
	}
	return newAccessToken, true
}

func validateHandler(w http.ResponseWriter, r *http.Request) bool {
	token := r.FormValue("accessToken")
	if validateAccessToken(token) {
		return true
	} else {
		return false
	}
}

func validateAccessToken(token string) bool {
	_, status := ParseJWT(token)
	return status == TokenValid
}
