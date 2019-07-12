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

/*
func (e *ViewEngine) NewInstance(name string) ViewRender {

	return ViewRender{
		Engine: e,
		Name:   name,
		Vars:   make(M),
	}
}
*/

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

func (r ViewRender) Render(w http.ResponseWriter) {
	err := r.Engine.executeRender(w, r.Name, r.Vars)
	if err != nil {
		switch e := err.(type) {
		case IStatusError:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}

}
