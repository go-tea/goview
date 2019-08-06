package goview

import (
	"log"
	"net/http"
)

// ViewRender struct
type ViewRender struct {
	Engine *ViewEngine
	Name   string
	Vars   M
}

// Instance method
func (e *ViewEngine) Instance(name string, data ...map[string]interface{}) ViewRender {

	d := make(M)

	if data != nil {
		d = data[0]
	}

	return ViewRender{
		Engine: e,
		Name:   name,
		Vars:   d,
	}
}

// Render method
func (r ViewRender) Render(w http.ResponseWriter) {
	err := r.Engine.executeRender(w, r.Name, r.Vars)
	if err != nil {
		switch t := err.(type) {
		case IStatusError:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", t.Status(), t)
			http.Error(w, t.Error(), t.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}

}
