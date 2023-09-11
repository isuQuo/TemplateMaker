package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/isuquo/templatemaker/internal/models"
	"github.com/isuquo/templatemaker/internal/rx"
	"github.com/isuquo/templatemaker/internal/validator"
)

type Result struct {
	Err     error
	Message string
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) respondWithJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// render is used to render a template to the client.
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", page))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "main", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

// decodePostForm decodes and validates the form data from a POST request.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

// isAuthenticated returns true if the current request is from an authenticated user.
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

// groupby groups templates by name.
func (app *application) groupby(userId string) map[string][]*models.Template {
	allTemplates, err := app.templates.SelectAll(userId)
	if err != nil {
		// TODO: Handle error
		panic(err)
	}

	groups := make(map[string][]*models.Template)
	for _, t := range allTemplates {
		groups[t.Name] = append(groups[t.Name], t)
	}

	return groups
}

func (app *application) importLog(file []byte) (map[string]interface{}, error) {
	var jsonObject map[string]interface{}
	if err := json.Unmarshal(file, &jsonObject); err != nil {
		return nil, err
	}

	return jsonObject, nil
}

// extractKeyValues extracts all the key-value pairs from a nested map into dot notation.
func (app *application) extractKeyValues(prefix string, data map[string]interface{}, kv *map[string]string) error {
	for key, value := range data {
		newPrefix := prefix
		if newPrefix != "" {
			newPrefix += "."
		}
		newPrefix += key

		switch v := value.(type) {
		case map[string]interface{}:
			app.extractKeyValues(newPrefix, v, kv)
		case []interface{}:
			// Create a list to store all items for combined representation
			var combinedItems []string

			for i, item := range v {
				switch itemValue := item.(type) {
				case map[string]interface{}:
					app.extractKeyValues(fmt.Sprintf("%s.%d", newPrefix, i), itemValue, kv)
				default:
					itemKey := fmt.Sprintf("%s.%d", newPrefix, i)
					itemStr := fmt.Sprintf("%v", itemValue)
					(*kv)[itemKey] = itemStr
					combinedItems = append(combinedItems, itemStr)
				}
			}

			// Set the combined representation of the list to the map
			(*kv)[newPrefix] = strings.Join(combinedItems, ",")
		default:
			(*kv)[newPrefix] = fmt.Sprintf("%v", value)
		}
	}

	return nil
}

// extractKeys extracts all the keys from a map.
func (app *application) extractKeys(kv map[string]string) []string {
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	return keys
}

func (app *application) getUserID(r *http.Request) string {
	return app.sessionManager.GetString(r.Context(), "authenticatedUserID")
}

func (app *application) processEmailTemplate(t *models.Template, files []*multipart.FileHeader) error {
	app.templates.UpdateStatus(t.ID, "in-progress")

	time.Sleep(5 * time.Second)

	_, err := rx.Test(t, files)
	if err != nil {
		app.templates.UpdateStatus(t.ID, "error")
		return err
	}

	app.templates.UpdateStatus(t.ID, "done")
	return nil
}

// Validate is used to validate the form data.
func (f *fileUploadForm) Validate() {
	f.CheckField(validator.IsJSONOrCSV(f.File), "file", "Invalid file content. Only JSON or CSV files are accepted.")
}

func (app *application) getStructs(t *models.Template) ([]rx.TestStruct, error) {
	//time.Sleep(5 * time.Second)

	return rx.GetStructs(t)
}
