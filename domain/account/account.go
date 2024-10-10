package account

type Account struct {
	ID     string `json:"id"      db:"id"`
	UserID string `json:"user_id" db:"user_id"`
	Income int    `json:"income"  db:"income"`
}
