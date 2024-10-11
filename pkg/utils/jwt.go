package utils

import (
	"financeapp/domain/user"
	"financeapp/pkg/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(u *user.User) (string, error) {
	claim := jwt.MapClaims{
		"iss":     "financeapp",
		"aud":     []string{"user"},
		"sub":     fmt.Sprintf("%s", u.Name),
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"user_id": u.ID,
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := claims.SignedString([]byte(config.Env.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
