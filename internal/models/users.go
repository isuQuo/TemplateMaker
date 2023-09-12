package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	IsAdmin        bool
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(id, name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (id, name, email, hashed_password, created, is_admin)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP(), FALSE)`

	_, err = m.DB.Exec(stmt, id, name, email, string(hashedPassword))
	if err != nil {
		// If this returns an error, we use errors.As() to check whether the error has the
		// underlying type *mysql.MySQLError. If it does, and the error number is 1062 (which
		// corresponds to a duplicate entry error) and the message contains the string
		// users_uc_email, then we know that the error relates to the users_uc_email unique
		// constraint on the email column. In this case, we return our own ErrDuplicateEmail
		// error.
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	var id string
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		} else {
			return "", err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		} else {
			return "", err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id string) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) Get(id string) (*User, error) {
	stmt := `SELECT id, name, email, created, is_admin FROM users WHERE id = ?`

	u := &User{}
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.IsAdmin)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (m *UserModel) SelectAll() ([]*User, error) {
	stmt := `SELECT id, name, email, created FROM users ORDER BY created DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		u := &User{}
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Created)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UserModel) Delete(id string) error {
	stmt := `DELETE FROM users WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
