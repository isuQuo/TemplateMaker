package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	JSONMimeType = "application/json"
	CSVMimeType  = "text/csv"
)

var AllowedFileExtensions = []string{".json", ".csv"}

// EmailRX is a regular expression that matches valid email addresses.
var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid() returns true if the FieldErrors map doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddNonFieldError() adds an error message to the NonFieldErrors slice
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// AddFieldError() adds an error mesage to the FieldErrors map if it doesn't exist
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField() adds an error message to the FieldErrors map
// only if a validation check is not 'ok'
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValues() returns true if a value is in a list of permitted values
func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}

	return false
}

// MinChars() returns true if a value contains at least n characters
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches() returns true if a value matches a regular expression
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// IsAllowedMimeType checks if the given value is a valid MIME type.
func IsAllowedMimeType(value string) bool {
	return value == JSONMimeType || value == CSVMimeType
}

// HasAllowedExtension checks if the given filename has an allowed extension.
func HasAllowedExtension(filename string) bool {
	for _, ext := range AllowedFileExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
