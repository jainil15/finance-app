package postgres

import (
	"database/sql"
	"errors"
	"financeapp/domain/user"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

type postgresUser struct {
	ID           uuid.UUID         `db:"id"`
	Name         user.Name         `db:"name"`
	Email        user.Email        `db:"email"`
	PasswordHash user.PasswordHash `db:"password_hash"`
}

func newPostgresUser(u *user.User) *postgresUser {
	return &postgresUser{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
}

func (u UserRepo) Add(us *user.User) (*user.User, error) {
	tx, err := u.db.Beginx()
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	dbUser := newPostgresUser(us)
	_, err = tx.NamedExec(
		"INSERT INTO users (id, name, email, password_hash) VALUES (:id, :name, :email, :password_hash);",
		dbUser,
	)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return nil, nil
}

func (u UserRepo) GetAll() ([]user.User, error) {
	var users []user.User
	rows, err := u.db.Queryx("SELECT * FROM users;")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		u := user.User{}
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (u UserRepo) GetById(userID uuid.UUID) (*user.User, error) {
	_user := postgresUser{}
	err := u.db.Get(
		&_user,
		"SELECT * FROM users WHERE id = $1",
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrorUserNotFound
		}
		return nil, err
	}
	us := user.New(_user.ID, _user.Name, _user.Email, _user.PasswordHash)
	return us, nil
}

func (u UserRepo) GetByEmail(email user.Email) (*user.User, error) {
	_user := postgresUser{}
	err := u.db.Get(
		&_user,
		"SELECT * FROM users WHERE email = $1",
		email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrorUserNotFound
		}
		return nil, err
	}
	us := user.New(_user.ID, _user.Name, _user.Email, _user.PasswordHash)
	return us, nil
}
