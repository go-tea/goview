//go:generate go-bindata  -pkg staticdata -o assets/static/bindata.go static/...
//go:generate go-bindata -prefix "views/" -pkg viewsdata -o assets/views/bindata.go views/...

package main

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/go-tea/goview"
	"github.com/go-tea/goview/_examples/bindata/assets/static"
	"github.com/go-tea/goview/_examples/bindata/assets/views"
	"github.com/go-tea/goview/assets/bindata"
	"github.com/go-tea/tea"
	"github.com/go-tea/tea/serve"
)

var e *goview.ViewEngine

func main() {

	router := tea.New(serve.RequestID, serve.RealIP, serve.Recoverer, serve.Logger)

	config := goview.Config{
		Root:         "views",
		Master:       "layouts/master",
		Extension:    ".html",
		DisableCache: true,
	}

	staticHandler := http.FileServer(&assetfs.AssetFS{Asset: staticdata.Asset, AssetDir: staticdata.AssetDir, AssetInfo: staticdata.AssetInfo, Prefix: "static"})
	router.GetSH("/static/*", http.StripPrefix("/static/", staticHandler))

	templateFS := &assetfs.AssetFS{Asset: viewsdata.Asset, AssetDir: viewsdata.AssetDir, AssetInfo: viewsdata.AssetInfo, Prefix: "views"}

	//new template engine
	e = bindata.NewWithConfig(templateFS, config)

	router.Get("/", h_home)
	router.Get("/page", h_page)

	// Start server
	router.ListenAndServe(":9090")
}

func h_home(w http.ResponseWriter, r *http.Request) {
	home := e.Instance("index", make(map[string]interface{}))
	home.Vars["title"] = "Index title!"
	home.Vars["add"] = func(a int, b int) int {
		return a + b
	}
	home.Render(w)

}

func h_page(w http.ResponseWriter, r *http.Request) {
	page := e.Instance("page.html", make(map[string]interface{}))
	page.Vars["title"] = "Page file title!!"
	page.Render(w)
}
