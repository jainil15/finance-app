package model

import (
	"financeapp/aggregate"
	"financeapp/domain/budget"
	"financeapp/domain/user"
)

type RegisterUser struct {
	Name     string
	Email    string
	Password string
}
type UserAggregate struct {
	user.User
	Transaction []aggregate.Transaction
	*budget.Budget
}

func NewUserAggregate(u *user.User, t *[]aggregate.Transaction, b *budget.Budget) *UserAggregate {
	return &UserAggregate{
		*u, *t, b,
	}
}

func NewRegisterUser(name, email, password string) RegisterUser {
	return RegisterUser{
		Name:     name,
		Email:    email,
		Password: password,
	}
}
