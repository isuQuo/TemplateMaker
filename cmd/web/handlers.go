package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/isuquo/templatemaker/internal/models"
	"github.com/isuquo/templatemaker/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type templateCreateForm struct {
	Name                string   `form:"name"`
	Subject             string   `form:"subject"`
	Description         string   `form:"description"`
	Assessment          []string `form:"assessment"`
	Recommendation      string   `form:"recommendation"`
	Query               []string `form:"query"`
	Status              string   `db:"status"`
	validator.Validator `form:"-"`
}

type userSignUpForm struct {
	Name                string     `form:"name"`
	Email               string     `form:"email"`
	Password            string     `form:"password"`
	validator.Validator `form:"-"` // This field is not a form field
}

type userSignInForm struct {
	Email               string     `form:"email"`
	Password            string     `form:"password"`
	validator.Validator `form:"-"` // This field is not a form field
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Templates = app.groupby(app.getUserID(r))

	app.render(w, http.StatusOK, "index.html", data)
}

func (app *application) templateCreateForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = templateCreateForm{}

	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) templateCreatePost(w http.ResponseWriter, r *http.Request) {
	var form templateCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// TODO: Implement error fields in form to render properly if errors occur.
	/*form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")
	form.CheckField(validator.PermittedValues(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")*/

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	// Generate a new, random ID for the template.
	id, err := uuid.NewUUID()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create a new Template struct containing the form data.
	template := &models.Template{
		ID:             id.String(),
		Name:           form.Name,
		Subject:        form.Subject,
		Description:    form.Description,
		Assessment:     strings.Join(form.Assessment, "{{EOA}}"),
		Recommendation: form.Recommendation,
		Query: sql.NullString{
			String: strings.Join(form.Query, "{{EOA}}"),
			Valid:  form.Query != nil,
		}.String,
		UserID: app.getUserID(r),
	}

	// Insert the template data in the database.
	err = app.templates.Insert(template)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the app.session.Put() method to add a message to the session.
	app.sessionManager.Put(r.Context(), "flash", "Template successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/template/view/%s", template.ID), http.StatusSeeOther)
}

func (app *application) templateEditForm(w http.ResponseWriter, r *http.Request) {
	// Use the httprouter.Param object to retrieve the value of the :id
	// parameter from the request URL.
	id, err := uuid.Parse(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	template, err := app.templates.Get(id.String())
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	assessments := strings.Split(template.Assessment, "{{EOA}}")
	queries := strings.Split(template.Query, "{{EOA}}")

	data := app.newTemplateData(r)
	data.Template = template
	data.Assessment = assessments
	data.Query = queries

	app.render(w, http.StatusOK, "edit.html", data)
}

func (app *application) templateEditPost(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var form templateCreateForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// TODO: Implement error fields in form to render properly if errors occur.
	/*form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")
	form.CheckField(validator.PermittedValues(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")*/

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "edit.html", data)
		return
	}

	// Create a new Template struct containing the form data.
	template := &models.Template{
		ID:             id.String(),
		Name:           form.Name,
		Subject:        form.Subject,
		Description:    form.Description,
		Assessment:     strings.Join(form.Assessment, "{{EOA}}"),
		Recommendation: form.Recommendation,
		Query: sql.NullString{
			String: strings.Join(form.Query, "{{EOA}}"),
			Valid:  form.Query != nil,
		}.String,
		UserID: app.getUserID(r),
	}

	// Update the template data in the database.
	err = app.templates.Update(template)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the app.session.Put() method to add a message to the session.
	app.sessionManager.Put(r.Context(), "flash", "Template successfully updated!")

	http.Redirect(w, r, fmt.Sprintf("/template/view/%s", template.ID), http.StatusSeeOther)
}

func (app *application) templateViewForm(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	template, err := app.templates.Get(id.String())
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	assessments := strings.Split(template.Assessment, "{{EOA}}")
	queries := strings.Split(template.Query, "{{EOA}}")

	data := app.newTemplateData(r)
	data.Template = template
	data.Assessment = assessments
	data.Query = queries

	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) templateDeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.templates.Delete(id.String())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the app.session.Put() method to add a message to the session.
	app.sessionManager.Put(r.Context(), "flash", "Template successfully deleted!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignUpForm{}

	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) signupUserPost(w http.ResponseWriter, r *http.Request) {
	var form userSignUpForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	// Generate a new, random ID for the new user.
	id, err := uuid.NewUUID()
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.users.Insert(id.String(), form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
}

func (app *application) signinUserForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignInForm{}

	app.render(w, http.StatusOK, "signin.html", data)
}

func (app *application) signinUserPost(w http.ResponseWriter, r *http.Request) {
	var form userSignInForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signin.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signin.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Renew the session token to prevent session fixation attacks.
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "flash", "You've been logged in successfully!")

	http.Redirect(w, r, "/template/create", http.StatusSeeOther)
}

func (app *application) signoutUser(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) split(w http.ResponseWriter, r *http.Request) {
	jsonObject, err := app.importLog("test.json")
	if err != nil {
		app.serverError(w, err)
		return
	}

	kv := make(map[string]string)
	err = app.extractKeyValues("", jsonObject, &kv)
	if err != nil {
		app.serverError(w, err)
		return
	}

	kvi, err := json.MarshalIndent(&kv, "", "\t")
	if err != nil {
		app.serverError(w, err)
		return
	}

	keys := app.extractKeys(kv)
	//specialKeys := logs.SpecialKeys

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"keys":     keys,
		"jsonData": string(kvi),
		//"special":  specialKeys,
	})
}

func (app *application) templateEmailPost(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	go app.processEmailTemplate(id.String())
	http.Redirect(w, r, fmt.Sprintf("/template/loading/%s", id), http.StatusSeeOther)
}

func (app *application) showLoading(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	template, err := app.templates.Get(id.String())
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	queries := strings.Split(template.Query, "{{EOA}}")

	data := app.newTemplateData(r)
	data.Template = template
	data.Query = queries

	app.render(w, http.StatusOK, "loading.html", data)
}

func (app *application) checkStatus(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")

	status, err := app.templates.GetStatus(id)
	if err != nil {
		app.serverError(w, err)
	}

	// Send the status back to the client.
	w.Write([]byte(status))
}

func (app *application) previewEmail(w http.ResponseWriter, r *http.Request) {
	// Write template to http writer temporarily
	w.Write([]byte("Previewing e-mail"))
}
