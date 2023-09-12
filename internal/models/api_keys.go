package models

import (
	"database/sql"
)

type APIKey struct {
	Name     string
	KeyValue string
}

type APIKeyModel struct {
	DB *sql.DB
}

func (m *APIKeyModel) Insert(name, apiKey string) error {
	stmt := `INSERT INTO api_keys (name, api_key)
	VALUES(?, ?)`

	_, err := m.DB.Exec(stmt, name, apiKey)
	if err != nil {
		return err
	}

	return nil
}

func (m *APIKeyModel) Get(name string) (*APIKey, error) {
	stmt := `SELECT name, api_key FROM api_keys WHERE name = ?`

	row := m.DB.QueryRow(stmt, name)

	var APIKey APIKey
	err := row.Scan(&APIKey.Name, &APIKey.KeyValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &APIKey, nil
}

func (m *APIKeyModel) GetAll() ([]*APIKey, error) {
	stmt := `SELECT name, api_key FROM api_keys ORDER BY name`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var APIKeys []*APIKey

	for rows.Next() {
		var APIKey APIKey
		err = rows.Scan(&APIKey.Name, &APIKey.KeyValue)
		if err != nil {
			return nil, err
		}

		APIKeys = append(APIKeys, &APIKey)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return APIKeys, nil
}

func (m *APIKeyModel) Update(name, apiKey string) error {
	stmt := `UPDATE api_keys SET name = ?, api_key = ? WHERE name = ?`

	_, err := m.DB.Exec(stmt, name, apiKey, name)
	if err != nil {
		return err
	}

	return nil
}

func (m *APIKeyModel) Delete(name string) error {
	stmt := `DELETE FROM api_keys WHERE name = ?`

	_, err := m.DB.Exec(stmt, name)
	if err != nil {
		return err
	}

	return nil
}
