package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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

type fileUploadForm struct {
	File                *multipart.FileHeader `form:"file"`
	validator.Validator `form:"-"`
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
	var form fileUploadForm

	// Extract the uploaded file from the request
	file, header, err := r.FormFile("logFile")
	if err != nil {
		app.respondWithJSONError(w, "Error reading file from request.", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	form.File = header
	form.Validate()

	if !form.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if !form.Valid() {
		app.respondWithJSONError(w, form.FieldErrors["file"], http.StatusBadRequest)
		return
	}

	// Read the file contents into a byte slice
	fileContent, err := io.ReadAll(file)
	if err != nil {
		app.serverError(w, fmt.Errorf("error reading file contents: %v", err))
		return
	}

	// Assuming app.importLog is modified to read from a byte slice
	// Replace this with your file parsing logic
	jsonObject, err := app.importLog(fileContent)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"keys":     keys,
		"jsonData": string(kvi),
	})
}

func (app *application) templateEmailPost(w http.ResponseWriter, r *http.Request) {
	var form fileUploadForm
	logFiles := []*multipart.FileHeader{}

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

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if r.MultipartForm != nil {
		logsRequired := len(strings.Split(template.Assessment, "{{EOA}}"))
		if len(r.MultipartForm.File) < logsRequired {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		for _, files := range r.MultipartForm.File {
			if len(files) > 0 {
				form.File = files[0]
				form.Validate()

				if !form.Valid() {
					app.respondWithJSONError(w, form.FieldErrors["file"], http.StatusBadRequest)
					return
				} else {
					// Add file to slice of files
					logFiles = append(logFiles, form.File)
				}
			}
		}
	}

	fmt.Println("Processing email template...")
	err = app.processEmailTemplate(template, logFiles)
	fmt.Println("Done processing email template...", err)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"status": "error", "message": "%s"}`, err.Error())))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/template/preview/%s", template.ID), http.StatusSeeOther)
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

func (app *application) getTemplateLogs(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")

	template, err := app.templates.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Send the status back to the client.
	logsRequired := len(strings.Split(template.Assessment, "{{EOA}}"))
	queries := strings.Split(template.Query, "{{EOA}}")

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"totalLogsRequired": logsRequired,
		"queries":           queries,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.serverError(w, err)
	}
}
