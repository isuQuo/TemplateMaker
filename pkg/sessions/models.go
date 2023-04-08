package sessions

import "github.com/isuquo/copper-test/pkg/users"

type Sessions struct {
	ID     string
	UserID string
	// Token is only set when the session is created
	Token     string
	TokenHash string
}

// Create creates a new session for the given user ID.
func (s *Sessions) Create(userId int) (*Sessions, error) {
	return nil, nil
}

// User returns the user associated with the given token.
func (s *Sessions) User(token string) (*users.User, error) {
	return nil, nil
}
