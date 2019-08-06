package ginview

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/go-tea/goview"
)

const templateEngineKey = "foolin-goview-ginview"

// ViewEngine struct
type ViewEngine struct {
	*goview.ViewEngine
}

// ViewRender struct
type ViewRender struct {
	Engine *ViewEngine
	Name   string
	Data   interface{}
}

// New function
func New(config goview.Config) *ViewEngine {
	return &ViewEngine{
		ViewEngine: goview.New(config),
	}
}

// Default function
func Default() *ViewEngine {
	return New(goview.DefaultConfig)
}

// Instance method
func (e *ViewEngine) Instance(name string, data interface{}) render.Render {
	return ViewRender{
		Engine: e,
		Name:   name,
		Data:   data,
	}
}

// HTML method
func (e *ViewEngine) HTML(ctx *gin.Context, code int, name string, data interface{}) {
	instance := e.Instance(name, data)
	ctx.Render(code, instance)
}

// Render (YAML) marshals the given interface object and writes data with custom ContentType.
func (v ViewRender) Render(w http.ResponseWriter) error {
	return v.Engine.RenderWriter(w, v.Name, v.Data)
}

// WriteContentType method
func (v ViewRender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = goview.HtmlContentType
	}
}

//NewMiddleware create new gin middleware
func NewMiddleware(config goview.Config) gin.HandlerFunc {
	return Middleware(New(config))
}

// Middleware function
// You should use helper func `Middleware()` to set the supplied
func Middleware(e *ViewEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(templateEngineKey, e)
	}
}

// HTML function
// TemplateEngine and make `HTML()` work validly.
func HTML(ctx *gin.Context, code int, name string, data interface{}) {
	if val, ok := ctx.Get(templateEngineKey); ok {
		if e, ok := val.(*ViewEngine); ok {
			e.HTML(ctx, code, name, data)
			return
		}
	}
	ctx.HTML(code, name, data)
}
