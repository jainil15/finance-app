package service

import (
	"financeapp/domain/account"
	"financeapp/domain/user"
	"financeapp/pkg/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserRegister struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type userResponse struct {
	ID    uuid.UUID  `json:"id"`
	Name  user.Name  `json:"name"`
	Email user.Email `json:"email"`
}
type accountResponse struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}
type registerResponse struct {
	UserResponse    *userResponse    `json:"user"`
	AccountResponse *accountResponse `json:"account"`
}

type UserService struct {
	userRepo    user.Repo
	accountRepo account.Repo
}

func NewRegisterResponse(ur *userResponse, ar *accountResponse) *registerResponse {
	return &registerResponse{
		ur, ar,
	}
}

func ToUserResponse(u *user.User) *userResponse {
	return &userResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func ToAccountResponse(a *account.Account) *accountResponse {
	return &accountResponse{
		ID:     a.ID,
		UserID: a.UserID,
	}
}

// NOTE: maybe move to routes package
func NewUserRoutes(e *echo.Echo, ur *UserService) {
	e.POST("/user/register", ur.Register)
	e.GET("/user", ur.GetAll)
}

func NewUserService(ur user.Repo, ar account.Repo) *UserService {
	return &UserService{
		userRepo:    ur,
		accountRepo: ar,
	}
}

func (u UserService) Register(c echo.Context) error {
	userRegister := UserRegister{}
	c.Bind(&userRegister)
	newUser, err := user.NewUser(userRegister.Name, userRegister.Email, userRegister.Password)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errs": fmt.Sprintf("%s", err),
		})
	}
	newAccount := account.New(newUser.ID)
	_, err = u.userRepo.Add(newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errs": fmt.Sprintf("%s", err),
		})
	}
	_, err = u.accountRepo.Add(newAccount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errs": fmt.Sprintf("%s", err),
		})
	}
	ar := ToAccountResponse(newAccount)
	ur := ToUserResponse(newUser)
	res := NewRegisterResponse(ur, ar)
	c.JSON(http.StatusOK, utils.Response{
		Message: "User Created",
		Result:  res,
	})
	return nil
}

func (u UserService) GetAll(c echo.Context) error {
	users, err := u.userRepo.GetAll()
	if err != nil {
		return err
	}
	uResponse := make([]*userResponse, len(users))
	for i, _user := range users {
		uResponse[i] = ToUserResponse(&_user)
	}
	log.Println(users)
	c.JSON(http.StatusOK, uResponse)
	return nil
}
