package users

import (
	"context"
	"database/sql"

	"github.com/gocopper/copper/csql"
)

var ErrRecordNotFound = sql.ErrNoRows

func NewQueries(querier csql.Querier) *Queries {
	return &Queries{
		querier: querier,
	}
}

type Queries struct {
	querier csql.Querier
}

func (q *Queries) CreateUser(ctx context.Context, user *User) error {
	const query = "INSERT INTO users (id, email, password_hash) VALUES (?, ?, ?)"

	_, err := q.querier.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	const query = "SELECT * from users where email=?"

	var (
		user User
		err  = q.querier.Get(ctx, &user, query, email)
	)

	return &user, err
}
