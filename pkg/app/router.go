package app

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
	"github.com/isuquo/copper-test/pkg/templates"
)

type NewRouterParams struct {
	RW        *chttp.ReaderWriter
	Logger    clogger.Logger
	Templates *templates.Queries
}

func NewRouter(p NewRouterParams) *Router {
	return &Router{
		rw:        p.RW,
		logger:    p.Logger,
		templates: p.Templates,
	}
}

type Router struct {
	rw        *chttp.ReaderWriter
	logger    clogger.Logger
	templates *templates.Queries
}

func (ro *Router) Routes() []chttp.Route {
	return []chttp.Route{
		{
			Path:    "/submit",
			Methods: []string{http.MethodPost},
			Handler: ro.HandleSubmitTemplate,
		},

		{
			Path:    "/submit",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleSubmitPage,
		},

		{
			Path:    "/",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleIndexPage,
		},
	}
}

func (ro *Router) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "index.html",
	})
}

func (ro *Router) HandleSubmitPage(w http.ResponseWriter, r *http.Request) {
	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "submit.html",
	})
}

func (ro *Router) HandleSubmitTemplate(w http.ResponseWriter, r *http.Request) {
	var (
		subject        = strings.TrimSpace(r.PostFormValue("title"))
		description    = strings.TrimSpace(r.PostFormValue("description"))
		assessment     = strings.TrimSpace(r.PostFormValue("assessment"))
		recommendation = strings.TrimSpace(r.PostFormValue("recommendation"))
	)

	if subject == "" || description == "" || assessment == "" || recommendation == "" {
		ro.rw.WriteHTMLError(w, r, cerrors.New(nil, "invalid template", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	err := ro.templates.SaveTemplate(r.Context(), &templates.Template{
		ID:             uuid.New().String(),
		Subject:        subject,
		Description:    description,
		Assessment:     "user1",
		Recommendation: "user2",
	})
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to save template", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
