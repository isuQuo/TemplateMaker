package models

import (
	"database/sql"
	"errors"
)

type Template struct {
	ID             string
	Name           string
	Subject        string
	Description    string
	Assessment     string
	Recommendation string
	Query          string
	Status         string
	UserID         string
}

type TemplateModel struct {
	DB *sql.DB
}

func (m *TemplateModel) SelectAll(userId string) ([]*Template, error) {
	query := `SELECT * FROM templates WHERE user_id = ?;`

	rows, err := m.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := []*Template{}

	for rows.Next() {
		s := &Template{}
		err = rows.Scan(&s.ID, &s.Name, &s.Subject, &s.Description, &s.Assessment, &s.Recommendation, &s.Query, &s.Status, &s.UserID)
		if err != nil {
			return nil, err
		}
		templates = append(templates, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return templates, nil
}

func (m *TemplateModel) Insert(template *Template) error {
	const query = `
	INSERT INTO templates (id, name, subject, description, assessment, recommendation, query, user_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	// Use the Exec() method on the embedded connection pool to execute our
	_, err := m.DB.Exec(query,
		template.ID,
		template.Name,
		template.Subject,
		template.Description,
		template.Assessment,
		template.Recommendation,
		template.Query,
		template.UserID,
	)

	return err
}

func (m *TemplateModel) Update(template *Template) error {
	const query = `
	UPDATE templates
	SET name=?, subject=?, description=?, assessment=?, recommendation=?, query=?
	WHERE id=?`

	_, err := m.DB.Exec(query,
		template.Name,
		template.Subject,
		template.Description,
		template.Assessment,
		template.Recommendation,
		template.Query,
		template.ID,
	)
	return err
}

func (m *TemplateModel) UpdateStatus(id string, status string) error {
	const query = `
	UPDATE templates
	SET status=?
	WHERE id=?`

	_, err := m.DB.Exec(query, status, id)
	return err
}

func (m *TemplateModel) GetStatus(id string) (string, error) {
	const query = `
	SELECT status
	FROM templates
	WHERE id=?`

	var status string
	err := m.DB.QueryRow(query, id).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (m *TemplateModel) Get(id string) (*Template, error) {
	const query = "SELECT * from templates where id=?"

	s := &Template{}

	err := m.DB.QueryRow(query, id).Scan(&s.ID, &s.Name, &s.Subject, &s.Description, &s.Assessment, &s.Recommendation, &s.Query, &s.Status, &s.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *TemplateModel) Delete(id string) error {
	const query = "DELETE FROM templates WHERE id=?"

	_, err := m.DB.Exec(query, id)
	return err
}
