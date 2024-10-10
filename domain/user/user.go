package user

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLength = 10
)

var (
	ErrorPasswordMismatch = errors.New("Incorrect Password")
	ErrorInvalidPassword  = errors.New("Password Not Valid")
	ErrorInvalidUUID      = errors.New("Invalid UUID")
)

type (
	PasswordHash string
	Email        string
	Name         string
	Password     string
)

type User struct {
	ID           uuid.UUID
	Name         Name
	Email        Email
	PasswordHash PasswordHash
}

func New(ID uuid.UUID, name Name, email Email, passwordHash PasswordHash) *User {
	return &User{
		ID, name, email, passwordHash,
	}
}

func NewUser(name string, email string, password string) (*User, error) {
	_id := uuid.New()
	_name, err := NewName(name)
	if err != nil {
		return nil, err
	}
	_email, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	_password, err := NewPassword(password)
	if err != nil {
		return nil, err
	}
	_passwordHash, err := NewPasswordHash(_password)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           _id,
		Name:         _name,
		Email:        _email,
		PasswordHash: _passwordHash,
	}, nil
}

func NewName(n string) (Name, error) {
	return Name(n), nil
}

func NewEmail(e string) (Email, error) {
	return Email(e), nil
}

func NewPassword(p string) (Password, error) {
	if len(p) < 10 {
		return Password(
				"",
			), errors.Join(
				ErrorInvalidPassword,
				errors.New("Password Length must be greater than 10"),
			)
	}
	return Password(p), nil
}

func NewPasswordHash(p Password) (PasswordHash, error) {
	pb, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		return "", err
	}
	return PasswordHash(string(pb)), err
}

func ComparePassword(p Password, ph PasswordHash) error {
	err := bcrypt.CompareHashAndPassword([]byte(ph), []byte(p))
	if err != nil {
		return ErrorInvalidPassword
	}
	return nil
}
