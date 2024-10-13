package middleware

import (
	"financeapp/pkg/config"
	"financeapp/pkg/utils"
	"financeapp/web/components/home"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessTokenCookie, err := c.Cookie("access-token")
		if err != nil {
			c.Response().Header().Add("HX-Replace-Url", "/login")
			return utils.WriteHTML(c, home.Home())
		}
		if err := accessTokenCookie.Valid(); err != nil {
			c.Response().Header().Add("HX-Replace-Url", "/login")
			log.Println("Error Validity:", err)
			return utils.WriteHTML(c, home.Home())
		}
		accessTokenString := accessTokenCookie.Value
		// accessTokenString := getAccessTokenFromRequest(c.Request())
		if accessTokenString == "" {
			c.Response().Header().Add("HX-Replace-Url", "/login")
			log.Println("Error Empty:", err)
			return utils.WriteHTML(c, home.Home())
		}
		accessToken, err := validateAccessToken(accessTokenString)
		if err != nil {
			c.Response().Header().Add("HX-Replace-Url", "/login")
			log.Println("Error Validity:", err)
			return utils.WriteHTML(c, home.Home())
		}
		if !accessToken.Valid {
			c.Response().Header().Add("HX-Replace-Url", "/login")
			log.Println("Error Invalid Again:", err)
			return utils.WriteHTML(c, home.Home())
		}

		claims, ok := accessToken.Claims.(jwt.MapClaims)
		if !ok {
			c.Response().Header().Add("HX-Replace-Url", "/login")
			log.Println("Error CLAIMS:", err)
			return utils.WriteHTML(c, home.Home())
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
	log.Printf("Token: %s\n", accessToken)
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
