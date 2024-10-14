package userService

import (
	"errors"
	"financeapp/domain/account"
	"financeapp/domain/budget"
	"financeapp/domain/category"
	"financeapp/domain/user"
	errx "financeapp/pkg/errors"
	"financeapp/pkg/middleware"
	"financeapp/pkg/model"
	"financeapp/pkg/utils"
	"financeapp/repository/repository"
	"financeapp/web/components/forms"
	"financeapp/web/layout"
	"financeapp/web/views"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserRegister struct {
	Name     string `json:"name"     form:"name"`
	Email    string `json:"email"    form:"email"`
	Password string `json:"password" form:"password"`
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
	userRepo        user.Repo
	accountRepo     account.Repo
	transactionRepo repository.TransactionRepo
	categoryRepo    category.Repo
	budgetRepo      budget.Repo
}
type LoginRequest struct {
	Email    string `json:"email"    form:"email"    validate:"email"`
	Password string `json:"password" form:"password" validate:"email"`
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

type loginResponse struct {
	UserResponse userResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
}

func NewLoginResponse(us *userResponse, token string) *loginResponse {
	return &loginResponse{
		UserResponse: *us,
		AccessToken:  token,
	}
}

func ToAccountResponse(a *account.Account) *accountResponse {
	return &accountResponse{
		ID:     a.ID,
		UserID: a.UserID,
	}
}

// NOTE: maybe move to routes package
func NewUserRoutes(g *echo.Group, ur *UserService) {
	g.POST("/user/register", ur.Register)
	g.GET("/user", middleware.AuthMiddleware(ur.GetAll))
	g.GET("/user/:user_id", middleware.AuthMiddleware(middleware.CheckUser(ur.GetById)))
	g.POST("/user/login", ur.Login)
}

func (us UserService) GetUserInfoView(c echo.Context) error {
	u := c.Get("user_id")
	log.Printf("USERID: %v\n", u)
	userIDString, ok := u.(string)
	if !ok {
		return c.JSON(400, utils.Error{
			Message: errors.New("Invalid UUID").Error(),
		})
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return err
	}
	userInfo, err := us.GetUserInfo(userID)
	if err != nil {
		return c.JSON(400, utils.Error{
			Message: errors.New("Invalid UUID").Error(),
		})
	}
	return utils.WriteHTML(c, layout.Layout(views.UserHome(*userInfo)))
}

func (us UserService) GetUserInfo(userID uuid.UUID) (*model.UserAggregate, error) {
	u, err := us.userRepo.GetById(userID)
	if err != nil {
		return nil, user.ErrorUserNotFound
	}

	t, err := us.transactionRepo.GetByUserId(userID)
	if err != nil {
		return nil, err
	}
	b, err := us.budgetRepo.GetByUserID(userID)
	if err != nil {
		if !errors.Is(err, budget.ErrorBudgetNotFound) {
			return nil, err
		}
		log.Println("Error Budget", err)
	}
	c, err := us.categoryRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	mUser := model.NewUserAggregate(*u, *t, c, b)
	fmt.Println(mUser)
	return mUser, nil
}

func NewUserService(
	ur user.Repo,
	ar account.Repo,
	tr repository.TransactionRepo,
	cr category.Repo,
	br budget.Repo,
) *UserService {
	return &UserService{
		userRepo:        ur,
		accountRepo:     ar,
		transactionRepo: tr,
		categoryRepo:    cr,
		budgetRepo:      br,
	}
}

func (lr *LoginRequest) ValidateUserLogin() errx.Error {
	errs := errx.New()
	if _, err := user.NewEmail(string(lr.Email)); err != nil {
		errs.Add("email", err.Error())
	}
	if len(lr.Password) == 0 {
		errs.Add("password", user.ErrorEmptyPassword.Error())
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (u UserService) Register(c echo.Context) error {
	userRegister := UserRegister{}
	err := c.Bind(&userRegister)
	if err != nil {
		errs := errx.New()
		errs.Add("server", err.Error())
		return utils.WriteHTML(
			c,
			forms.Register(
				model.NewRegisterUser(userRegister.Name, userRegister.Email, userRegister.Password),
				errs,
			),
		)
	}
	newUser, err := user.NewUser(userRegister.Name, userRegister.Email, userRegister.Password)
	if err != nil {
		errs := errx.New()
		errs.Add("server", err.Error())
		return utils.WriteHTML(
			c,
			forms.Register(
				model.NewRegisterUser(userRegister.Name, userRegister.Email, userRegister.Password),
				errs,
			),
		)
	}
	newAccount := account.New(newUser.ID)
	_, err = u.userRepo.Add(newUser)
	if err != nil {
		errs := errx.New()
		errs.Add("server", err.Error())
		return utils.WriteHTML(
			c,
			forms.Register(
				model.NewRegisterUser(userRegister.Name, userRegister.Email, userRegister.Password),
				errs,
			),
		)
	}
	_, err = u.accountRepo.Add(newAccount)
	if err != nil {
		errs := errx.New()
		errs.Add("server", err.Error())
		return utils.WriteHTML(
			c,
			forms.Register(
				model.NewRegisterUser(userRegister.Name, userRegister.Email, userRegister.Password),
				errs,
			),
		)
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
	log.Println(c.Get("user_id"))
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: fmt.Sprintf("Invalid user id: %s", c.Param("user_id")),
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
		errs := errx.New()
		errs.Add("server", err.Error())
		return utils.WriteHTML(
			c,
			forms.Login(loginUser.Email, loginUser.Password, &errs),
		)
	}
	errs := loginUser.ValidateUserLogin()
	if errs != nil {
		return utils.WriteHTML(
			c,
			forms.Login(loginUser.Email, loginUser.Password, &errs),
		)
	}
	us, err := u.userRepo.GetByEmail(user.Email(loginUser.Email))
	if err != nil {
		if errors.Is(err, user.ErrorUserNotFound) {
			errs = errx.New()
			errs.Add("email", fmt.Sprintf("User with user email: %s not found", loginUser.Email))
			return utils.WriteHTML(
				c,
				forms.Login(loginUser.Email, loginUser.Password, &errs),
			)
		}

		return c.JSON(http.StatusInternalServerError, utils.Error{
			Message: err.Error(),
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(us.PasswordHash), []byte(loginUser.Password))
	if err != nil {
		if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) {
			errs = errx.New()
			errs.Add("password", "Invalid Password")
			return utils.WriteHTML(
				c,
				forms.Login(loginUser.Email, loginUser.Password, &errs),
			)
		}
		errs = errx.New()
		errs.Add("server", "Internal server error")
		return utils.WriteHTML(
			c,
			forms.Login(loginUser.Email, loginUser.Password, &errs),
		)
	}
	token, err := utils.CreateToken(us)
	if err != nil {
		errs = errx.New()
		errs.Add("server", "Internal server error")
		return utils.WriteHTML(
			c,
			forms.Login(loginUser.Email, loginUser.Password, &errs),
		)
	}
	accessTokenCookie := http.Cookie{
		Name:     "access-token",
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   8000000,
		Path:     "/",
		Value:    token,
	}
	c.SetCookie(&accessTokenCookie)
	usInfo, err := u.GetUserInfo(us.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	c.Response().Header().Set("HX-Replace-Url", "/home")
	return utils.WriteHTML(c, views.UserHome(*usInfo))
}
