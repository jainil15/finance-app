package service

import (
	"financeapp/domain/user"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserRegister struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService struct {
	userRepo user.Repo
}

// NOTE: maybe move to routes package
func NewUserRoutes(e *echo.Echo, ur *UserService) {
	e.Logger.Debugf("Hello")
	e.POST("/user/register", ur.Register)
}

func NewUserService(ur user.Repo) *UserService {
	return &UserService{
		userRepo: ur,
	}
}

func (u UserService) Register(c echo.Context) error {
	ur := UserRegister{}
	c.Bind(&ur)
	user, err := user.NewUser(ur.Name, ur.Email, ur.Password)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errs": fmt.Sprintf("%s", err),
		})
	}
	u.userRepo.Add(user)
	c.JSON(http.StatusOK, user)
	return nil
}
