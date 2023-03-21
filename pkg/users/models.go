package users

type User struct {
	ID           string
	Email        string
	PasswordHash string `db:"password_hash"`
}
