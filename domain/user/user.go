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
	ErrorEmptyName        = errors.New("Name required")
	ErrorPasswordMismatch = errors.New("Incorrect Password")
	ErrorEmptyEmail       = errors.New("Email Required")
	ErrorInvalidPassword  = errors.New("Password Not Valid")
	ErrorInvalidUUID      = errors.New("Invalid UUID")
	ErrorPasswordTooShort = errors.New("Password too short")
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
	if len(n) < 1 {
		return "", ErrorEmptyEmail
	}
	return Name(n), nil
}

func NewEmail(e string) (Email, error) {
	if len(e) < 1 {
		return "", ErrorEmptyEmail
	}
	return Email(e), nil
}

func NewPassword(p string) (Password, error) {
	if len(p) < minPasswordLength {
		return Password(""), ErrorPasswordTooShort
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
