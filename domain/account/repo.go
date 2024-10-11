package account

type Repo interface {
	Add(account *Account) (*Account, error)
}
