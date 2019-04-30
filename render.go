package goview

import (
	"log"
	"net/http"
)

type ViewRender struct {
	Engine *ViewEngine
	Name   string
	Vars   M
}

func (e *ViewEngine) Instance(name string, data map[string]interface{}) ViewRender {

	return ViewRender{
		Engine: e,
		Name:   name,
		Vars:   data,
	}
}

func (r ViewRender) Render(w http.ResponseWriter) {
	err := r.Engine.executeRender(w, r.Name, r.Vars)
	if err != nil {
		switch e := err.(type) {
		case IStatusError:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			if e.Error != nil {
				log.Printf("HTTP %d - %s", e.Status(), e)
				http.Error(w, e.Error(), e.Status())
			}
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}

}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type IStatusError interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}
