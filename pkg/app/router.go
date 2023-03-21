package app

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
	"github.com/isuquo/copper-test/pkg/logs"
	"github.com/isuquo/copper-test/pkg/templates"
	"github.com/isuquo/copper-test/pkg/users"
)

type NewRouterParams struct {
	RW        *chttp.ReaderWriter
	Logger    clogger.Logger
	Templates *templates.Queries
	Users     *users.Queries
}

func NewRouter(p NewRouterParams) *Router {
	return &Router{
		rw:        p.RW,
		logger:    p.Logger,
		templates: p.Templates,
		users:     p.Users,
	}
}

type Router struct {
	rw        *chttp.ReaderWriter
	logger    clogger.Logger
	templates *templates.Queries
	users     *users.Queries
}

func (ro *Router) Routes() []chttp.Route {
	return []chttp.Route{
		{
			Path:    "/signin",
			Methods: []string{http.MethodPost},
			Handler: ro.HandleUserSignin,
		},

		{
			Path:    "/signin",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleUserSigninPage,
		},

		{
			Path:    "/signup",
			Methods: []string{http.MethodPost},
			Handler: ro.HandleUserSignup,
		},

		{
			Path:    "/signup",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleUserSignupPage,
		},

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
	type TemplateData struct {
		Template templates.Template
		Log      string
		Keys     []string
	}

	jsonObject, err := logs.ImportJSONFile("test.json")
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to import log", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	kv := make(map[string]string)
	err = logs.ExtractKeyValues("", jsonObject, &kv)
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to extract log", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	kvi, err := json.MarshalIndent(&kv, "", "\t")
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to marshal log", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	keys := logs.ExtractKeys(kv)

	data := TemplateData{
		Log:  string(kvi),
		Keys: keys,
	}

	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "submit-split.html",
		Data:         data,
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
		Log      string
		Keys     []string
	}

	jsonObject, err := logs.ImportJSONFile("test.json")
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to import log", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	kv := make(map[string]string)
	err = logs.ExtractKeyValues("", jsonObject, &kv)
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to extract log", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	kvi, err := json.MarshalIndent(&kv, "", "\t")
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to marshal log", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	keys := logs.ExtractKeys(kv)

	data := TemplateData{
		Template: *template,
		Log:      string(kvi),
		Keys:     keys,
	}

	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "edit-split.html",
		Data:         data,
	})
}

func (ro *Router) HandleUserSignupPage(w http.ResponseWriter, r *http.Request) {
	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "signup.html",
	})
}

func (ro *Router) HandleUserSignup(w http.ResponseWriter, r *http.Request) {
	var (
		id       = uuid.New().String()
		email    = strings.TrimSpace(r.PostFormValue("email"))
		password = strings.TrimSpace(r.PostFormValue("password"))
	)

	if email == "" || password == "" {
		ro.rw.WriteHTMLError(w, r, cerrors.New(nil, "invalid signup details", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to hash password", map[string]interface{}{
			"form": r.Form,
		}))
	}
	passwordHash := string(hashedBytes)

	err = ro.users.CreateUser(r.Context(), &users.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to create user", map[string]interface{}{
			"form": r.Form,
		}))
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ro *Router) HandleUserSigninPage(w http.ResponseWriter, r *http.Request) {
	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "signin.html",
	})
}

func (ro *Router) HandleUserSignin(w http.ResponseWriter, r *http.Request) {
	var (
		email    = strings.TrimSpace(r.PostFormValue("email"))
		password = strings.TrimSpace(r.PostFormValue("password"))
	)

	if email == "" || password == "" {
		ro.rw.WriteHTMLError(w, r, cerrors.New(nil, "invalid signin details", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	user, err := ro.users.GetUserByEmail(r.Context(), email)
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to get user", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		ro.rw.WriteHTMLError(w, r, cerrors.New(err, "failed to compare password", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
