package user

type Repo interface {
	Add(user *User) (*User, error)
	GetAll() ([]*User, error)
}
