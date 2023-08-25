package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")

// ErrInvalidCredentials is returned when a user tries to authenticate with an
// invalid email address or password.
var ErrInvalidCredentials = errors.New("models: invalid credentials")

// ErrDuplicateEmail is returned when we attempt to create a user with an email
// address that is already in use.
var ErrDuplicateEmail = errors.New("models: duplicate email")
