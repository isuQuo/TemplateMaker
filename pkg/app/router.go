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
			Path:    "/edit-split/{id}",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleEditSplitPage,
		},

		{
			Path:    "/submit-split",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleSubmitSplitPage,
		},

		{
			Path:    "/edit/{id}",
			Methods: []string{http.MethodPost},
			Handler: ro.HandleEditTemplate,
		},

		{
			Path:    "/edit/{id}",
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
	id := string(chttp.URLParams(r)["id"])
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

// TODO: Parse Log and write JSON object to writer
func (ro *Router) HandleSubmitSplitPage(w http.ResponseWriter, r *http.Request) {
	var jsonObject = `{
	"Log": "RETURNING LOG"
}`

	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "submit-split.html",
		Data: map[string]interface{}{
			"Log": jsonObject,
		},
	})
}

func (ro *Router) HandleEditSplitPage(w http.ResponseWriter, r *http.Request) {
	id := string(chttp.URLParams(r)["id"])
	template, err := ro.templates.GetTemplateByID(r.Context(), id)
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to get template", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	type TemplateData struct {
		Template templates.Template
		Log      map[string]interface{}
	}

	// Replace with file object
	var jsonObject = `{
		"name": "John Doe",
		"age": 30,
		"address": {
		  "street": "123 Main St",
		  "city": "Anytown",
		  "state": "CA",
		  "zip": "12345"
		},
		"phone_numbers": [
		  {
			"type": "home",
			"number": "555-1234"
		  },
		  {
			"type": "work",
			"number": "555-5678"
		  }
		]
	  }`

	data := TemplateData{
		Template: *template,
		Log:      map[string]interface{}{"Log": jsonObject},
	}

	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "edit-split.html",
		Data:         data,
	})
}
