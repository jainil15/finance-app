package middleware

import (
	"financeapp/pkg/config"
	"financeapp/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessTokenString := getAccessTokenFromRequest(c.Request())
		if accessTokenString == "" {
			return c.JSON(http.StatusBadRequest, utils.Error{
				Message: "Missing Access Token",
			})
		}
		accessToken, err := validateAccessToken(accessTokenString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.Error{
				Message: "Access Token Malformed or Expired",
			})
		}
		if !accessToken.Valid {
			return c.JSON(http.StatusUnauthorized, utils.Error{
				Message: "Invalid Token",
			})
		}

		claims, ok := accessToken.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, utils.Error{
				Message: "Invalid Claims",
			})
		}
		log.Println(claims)
		c.Set("user_id", claims["user_id"])
		return next(c)
	}
}

func getAccessTokenFromRequest(r *http.Request) string {
	accessToken := r.Header.Get("Authorization")
	accessTokenSplit := strings.Split(strings.TrimSpace(accessToken), " ")
	if len(accessTokenSplit) < 2 {
		return ""
	}
	return strings.TrimSpace(strings.Split(accessToken, " ")[1])
}

func validateAccessToken(accessToken string) (*jwt.Token, error) {
	return jwt.Parse(
		accessToken,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected error")
			}
			return []byte(config.Env.JWTSecret), nil
		},
	)
}

func CheckUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ruserID := c.Param("user_id")
		userID := c.Get("user_id")
		if userID == "" {
			return c.JSON(http.StatusUnauthorized, utils.Error{
				Message: "User unauthorized",
			})
		}
		if ruserID != userID {
			return c.JSON(http.StatusForbidden, utils.Error{
				Message: "Access Forbidden",
			})
		}
		return next(c)
	}
}
