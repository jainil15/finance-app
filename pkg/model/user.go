package model

import (
	"financeapp/aggregate"
	"financeapp/domain/budget"
	"financeapp/domain/category"
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
	Category    []category.Category
	*budget.Budget
}

func NewUserAggregate(
	u user.User,
	t []aggregate.Transaction,
	c []category.Category,
	b *budget.Budget,
) *UserAggregate {
	return &UserAggregate{
		u, t, c, b,
	}
}

func NewRegisterUser(name, email, password string) RegisterUser {
	return RegisterUser{
		Name:     name,
		Email:    email,
		Password: password,
	}
}
