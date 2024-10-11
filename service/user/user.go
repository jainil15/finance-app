package service

import (
	"errors"
	"financeapp/domain/account"
	"financeapp/domain/user"
	"financeapp/pkg/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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
type LoginRequest struct {
	Email    user.Email    `json:"email"`
	Password user.Password `json:"password"`
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
	e.GET("/user/:user_id", ur.GetById)
	e.POST("/user/login", ur.Login)
}

func NewUserService(ur user.Repo, ar account.Repo) *UserService {
	return &UserService{
		userRepo:    ur,
		accountRepo: ar,
	}
}

func (lr *LoginRequest) ValidateUserLogin() {
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

func (u UserService) GetById(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: fmt.Sprintf("Invalid user id: %s", userID),
		})
	}
	us, err := u.userRepo.GetById(userID)
	if err != nil {
		if errors.Is(err, user.ErrorUserNotFound) {
			return c.JSON(http.StatusNotFound, utils.Error{
				Message: fmt.Sprintf("User with user id: %s not found", userID),
			})
		}
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Message: err.Error(),
		})
	}
	return c.JSON(
		http.StatusOK,
		ToUserResponse(us),
	)
}

func (u UserService) Login(c echo.Context) error {
	loginUser := LoginRequest{}
	err := c.Bind(&loginUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: fmt.Sprintf("Invalid request payload"),
			Error:   err,
		})
	}
	us, err := u.userRepo.GetByEmail(loginUser.Email)
	if err != nil {
		if errors.Is(err, user.ErrorUserNotFound) {
			return c.JSON(http.StatusNotFound, utils.Error{
				Message: fmt.Sprintf("User with user email: %s not found", loginUser.Email),
			})
		}
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Message: err.Error(),
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(us.PasswordHash), []byte(loginUser.Password))
	if err != nil {
		if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) {
			return c.JSON(http.StatusInternalServerError, utils.Error{
				Message: "Invalid Email or Password",
			})
		}
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Message: err.Error(),
		})
	}
	return c.JSON(
		http.StatusOK,
		ToUserResponse(us),
	)
}
