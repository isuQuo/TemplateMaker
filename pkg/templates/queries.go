package templates

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

func (q *Queries) ListTemplates(ctx context.Context) ([]Template, error) {
	const query = "SELECT * FROM templates;"

	var (
		templates []Template
		err       = q.querier.Select(ctx, &templates, query)
	)

	return templates, err
}

func (q *Queries) GetTemplateByID(ctx context.Context, id string) (*Template, error) {
	const query = "SELECT * from templates where id=?"

	var (
		template Template
		err      = q.querier.Get(ctx, &template, query, id)
	)

	return &template, err
}

func (q *Queries) SaveTemplate(ctx context.Context, template *Template) error {
	const query = `
	INSERT INTO templates (id, name, subject, description, assessment, recommendation)
	VALUES (?, ?, ?, ?, ?, ?)`
	//ON CONFLICT (id) DO UPDATE SET name=?, subject=?, description=?, assessment=?, recommendation=?`

	//assessment := strings.Split(template.Assessment, "{{EOA}}")
	_, err := q.querier.Exec(ctx, query,
		template.ID,
		template.Name,
		template.Subject,
		template.Description,
		template.Assessment,
		template.Recommendation,
	)

	return err
}

func (q *Queries) EditTemplate(ctx context.Context, template *Template) error {
	const query = `
	UPDATE templates
	SET name=?, subject=?, description=?, assessment=?, recommendation=?
	WHERE id=?`

	_, err := q.querier.Exec(ctx, query,
		template.Name,
		template.Subject,
		template.Description,
		template.Assessment,
		template.Recommendation,
		template.ID,
	)

	return err
}
