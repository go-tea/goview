package renderer

//github.com/thedevsaddam/renderer

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	// ContentType represents content type
	ContentType = "Content-Type"
	// ContentBinary header value for binary data.
	ContentBinary = "application/octet-stream"
	// ContentHTML header value for HTML data.
	ContentHTML = "text/html"
	// ContentJSON header value for JSON data.
	ContentJSON = "application/json"
	// ContentJSONP header value for JSONP data.
	ContentJSONP = "application/javascript"
	// ContentLength header constant.
	ContentLength = "Content-Length"
	// ContentText header value for Text data.
	ContentText = "text/plain"
	// ContentXHTML header value for XHTML data.
	ContentXHTML = "application/xhtml+xml"
	// ContentXML header value for XML data.
	ContentXML = "text/xml"
	// ContentYAML represents content type application/x-yaml
	ContentYAML = "application/x-yaml"
	// ContentOctet describes octet-stream
	ContentOctet = "application/octet-stream"
	// ContentDisposition describes contentDisposition
	ContentDisposition = "Content-Disposition"
	// contentDispositionInline describes content disposition type
	contentDispositionInline string = "inline"
	// contentDispositionAttachment describes content disposition type
	contentDispositionAttachment string = "attachment"
)

// Vars
var (
	JSONPrefix   string
	JSONIndent   bool
	XMLPrefix    string
	XMLIndent    bool
	UnEscapeHTML bool
)

// NoContent serve success but no content response
func NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Raw render serve raw response where you have to build the headers, body
func Raw(w http.ResponseWriter, status int, v interface{}) error {
	w.WriteHeader(status)
	_, err := w.Write(v.([]byte))
	return err
}

// String serve string content as text/plain response
func String(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, ContentText)
	w.WriteHeader(status)
	_, err := w.Write([]byte(v.(string)))
	return err
}

// JSON serve data as JSON as response
func JSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, ContentJSON)
	w.WriteHeader(status)

	bs, err := fjson(v)
	if err != nil {
		return err
	}
	if JSONPrefix != "" {
		w.Write([]byte(JSONPrefix))
	}
	_, err = w.Write(bs)
	return err
}

// JSONP serve data as JSONP response
func JSONP(w http.ResponseWriter, status int, callback string, v interface{}) error {
	w.Header().Set(ContentType, ContentJSONP)
	w.WriteHeader(status)

	bs, err := fjson(v)
	if err != nil {
		return err
	}

	if callback == "" {
		return errors.New("renderer: callback can not bet empty")
	}

	w.Write([]byte(callback + "("))
	_, err = w.Write(bs)
	w.Write([]byte(");"))

	return err
}

// XML serve data as XML response
func XML(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, ContentXML)
	w.WriteHeader(status)
	var bs []byte
	var err error

	if XMLIndent {
		bs, err = xml.MarshalIndent(v, "", " ")
	} else {
		bs, err = xml.Marshal(v)
	}
	if err != nil {
		return err
	}

	if XMLPrefix != "" {
		w.Write([]byte(XMLPrefix))
	}
	_, err = w.Write(bs)
	return err
}

// YAML serve data as YAML response
func YAML(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, ContentYAML)
	w.WriteHeader(status)

	bs, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}

// Binary serve file as application/octet-stream response; you may add ContentDisposition by your own.
func Binary(w http.ResponseWriter, status int, reader io.Reader, filename string, inline bool) error {
	if inline {
		w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", contentDispositionInline, filename))
	} else {
		w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", contentDispositionAttachment, filename))
	}
	w.Header().Set(ContentType, ContentBinary)
	w.WriteHeader(status)
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	_, err = w.Write(bs)
	return err
}

// File serve file as response from io.Reader
func File(w http.ResponseWriter, status int, reader io.Reader, filename string, inline bool) error {
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	// set headers
	mime := http.DetectContentType(bs)
	if inline {
		w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", contentDispositionInline, filename))
	} else {
		w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", contentDispositionAttachment, filename))
	}
	w.Header().Set(ContentType, mime)
	w.WriteHeader(status)

	_, err = w.Write(bs)
	return err
}

// FileView serve file as response with content-disposition value inline
func FileView(w http.ResponseWriter, status int, fpath, name string) error {
	return file(w, status, fpath, name, contentDispositionInline)
}

// FileDownload serve file as response with content-disposition value attachment
func FileDownload(w http.ResponseWriter, status int, fpath, name string) error {
	return file(w, status, fpath, name, contentDispositionAttachment)
}

// HTMLString render string as html. Note: You must provide trusted html when using this method
func HTMLString(w http.ResponseWriter, status int, html string) error {
	w.Header().Set(ContentType, ContentHTML)
	w.WriteHeader(status)
	out := template.HTML(html)
	_, err := w.Write([]byte(out))
	return err
}

/*
func MsgPack(w http.ResponseWriter, obj interface{}) error {
	w.Header().Set(ContentType, ContentMsgPack)
	.....
}
*/

// json converts the data as bytes using json encoder
func fjson(v interface{}) ([]byte, error) {
	var bs []byte
	var err error
	if JSONIndent {
		bs, err = json.MarshalIndent(v, "", " ")
	} else {
		bs, err = json.Marshal(v)
	}
	if err != nil {
		return bs, err
	}
	if UnEscapeHTML {
		bs = bytes.Replace(bs, []byte("\\u003c"), []byte("<"), -1)
		bs = bytes.Replace(bs, []byte("\\u003e"), []byte(">"), -1)
		bs = bytes.Replace(bs, []byte("\\u0026"), []byte("&"), -1)
	}
	return bs, nil
}

// file serve file as response
func file(w http.ResponseWriter, status int, fpath, name, contentDisposition string) error {
	var bs []byte
	var err error
	bs, err = ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(bs)

	// filename, ext, mimes
	var fn, mime, ext string
	fn, err = filepath.Abs(fpath)
	if err != nil {
		return err
	}
	ext = filepath.Ext(fpath)
	if name != "" {
		if !strings.HasSuffix(name, ext) {
			fn = name + ext
		}
	}

	mime = http.DetectContentType(bs)

	// set headers
	w.Header().Set(ContentType, mime)
	w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", contentDisposition, fn))
	w.WriteHeader(status)

	if _, err = buf.WriteTo(w); err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}
