package web

import (
	"embed"
	"net/http"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/chttp"
	"github.com/isuquo/copper-test/pkg/templates"
)

//go:embed src
var HTMLDir embed.FS

func HTMLRenderFuncs(q *templates.Queries) []chttp.HTMLRenderFunc {
	return []chttp.HTMLRenderFunc{
		{
			Name: "groupby",
			Func: func(r *http.Request) any {
				return func(map[string][]templates.Template) (map[string][]templates.Template, error) {
					allTemplates, err := q.ListTemplates(r.Context())
					if err != nil {
						return nil, cerrors.New(err, "failed to list templates", map[string]interface{}{
							"request": err,
						})
					}

					groups := make(map[string][]templates.Template)
					for _, t := range allTemplates {
						groups[t.Name] = append(groups[t.Name], t)
					}
					return groups, nil
				}
			},
		},
	}
}
