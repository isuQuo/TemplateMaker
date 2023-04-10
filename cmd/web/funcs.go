package main

import (
	"encoding/json"
	"html/template"
	"time"
)

var functions = template.FuncMap{
	"humanDate": humanDate,
	"json":      renderJson,
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func renderJson(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		// TODO: log error
		panic(err)
	}
	return string(b), nil
}
