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
			Path:    "/edit",
			Methods: []string{http.MethodPut},
			Handler: ro.HandleEditTemplate,
		},

		{
			Path:    "/edit",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleEditPage,
		},

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
	allTemplates, err := ro.templates.ListTemplates(r.Context())
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to list templates", nil))
		return
	}

	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "index.html",
		Data: map[string][]templates.Template{
			"Templates": allTemplates,
		},
	})
}

func (ro *Router) HandleSubmitPage(w http.ResponseWriter, r *http.Request) {
	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "submit.html",
	})
}

func (ro *Router) HandleSubmitTemplate(w http.ResponseWriter, r *http.Request) {
	var (
		name           = strings.TrimSpace(r.PostFormValue("name"))
		subject        = strings.TrimSpace(r.PostFormValue("subject"))
		description    = strings.TrimSpace(r.PostFormValue("description"))
		assessment     = strings.TrimSpace(r.PostFormValue("assessment"))
		recommendation = strings.TrimSpace(r.PostFormValue("recommendation"))
	)

	if name == "" || subject == "" || description == "" || assessment == "" || recommendation == "" {
		ro.rw.WriteHTMLError(w, r, cerrors.New(nil, "invalid template", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	err := ro.templates.SaveTemplate(r.Context(), &templates.Template{
		ID:             uuid.New().String(),
		Name:           name,
		Subject:        subject,
		Description:    description,
		Assessment:     assessment,
		Recommendation: recommendation,
	})
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to save template", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ro *Router) HandleEditPage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	template, err := ro.templates.GetTemplateByID(r.Context(), id)
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to get template", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "edit.html",
		Data:         template,
	})
}

func (ro *Router) HandleEditTemplate(w http.ResponseWriter, r *http.Request) {
	var (
		id             = string(chttp.URLParams(r)["id"])
		name           = strings.TrimSpace(r.PostFormValue("name"))
		subject        = strings.TrimSpace(r.PostFormValue("subject"))
		description    = strings.TrimSpace(r.PostFormValue("description"))
		assessment     = strings.TrimSpace(r.PostFormValue("assessment"))
		recommendation = strings.TrimSpace(r.PostFormValue("recommendation"))
	)

	if name == "" || subject == "" || description == "" || assessment == "" || recommendation == "" {
		ro.rw.WriteHTMLError(w, r, cerrors.New(nil, "invalid template", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	err := ro.templates.EditTemplate(r.Context(), &templates.Template{
		ID:             id,
		Name:           name,
		Subject:        subject,
		Description:    description,
		Assessment:     assessment,
		Recommendation: recommendation,
	})
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to edit template", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
