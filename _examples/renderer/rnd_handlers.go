package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"

	"github.com/go-tea/goview/renderer"
)

func home(w http.ResponseWriter, r *http.Request) {
	renderer.HTMLString(w, http.StatusOK, "<a href=\"/text\">text</a></br><a href=\"/xml\">xml</a></br>")
}

func nocontent(w http.ResponseWriter, r *http.Request) {
	renderer.NoContent(w)
}

func text(w http.ResponseWriter, r *http.Request) {
	renderer.String(w, http.StatusOK, "plain text")
}

type ExampleXml struct {
	XMLName xml.Name `xml:"example"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

func lmx(w http.ResponseWriter, r *http.Request) {
	renderer.XML(w, http.StatusOK, ExampleXml{One: "hello", Two: "xml"})
}

func json(w http.ResponseWriter, r *http.Request) {
	renderer.JSON(w, http.StatusOK, map[string]string{"hello": "json"})
}

func jsonp(w http.ResponseWriter, r *http.Request) {
	usr := struct {
		Name string
		Age  int
	}{"John Doe", 30}

	renderer.JSONP(w, http.StatusOK, "callback", usr)
}

func yaml(w http.ResponseWriter, r *http.Request) {
	usr := struct {
		Name string
		Age  int
	}{"John Doe", 30}
	renderer.YAML(w, http.StatusOK, usr)
}

func binary(w http.ResponseWriter, r *http.Request) {
	var reader io.Reader
	reader, _ = os.Open("main.go")
	renderer.Binary(w, http.StatusOK, reader, "main.go", true)
}

func fileinline(w http.ResponseWriter, r *http.Request) {
	renderer.FileView(w, http.StatusOK, "main.go", "main.go")
}

func filedownload(w http.ResponseWriter, r *http.Request) {
	renderer.FileDownload(w, http.StatusOK, "main.go", "main.go")
}
