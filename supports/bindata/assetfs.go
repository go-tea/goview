package bindata

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/go-tea/goview"
)

// go-bindata -prefix "dir/" -pkg main -o bindata.go dir/...

/**
New template engine, default views root.
*/
func New(viewsRootBox *assetfs.AssetFS) *goview.ViewEngine {
	return NewWithConfig(viewsRootBox, goview.DefaultConfig)
}

func NewWithConfig(viewsRootBox *assetfs.AssetFS, config goview.Config) *goview.ViewEngine {
	config.Root = viewsRootBox.Prefix
	engine := goview.New(config)
	engine.SetFileHandler(FileHandler(viewsRootBox))
	return engine
}

func FileHandler(viewsRootBox *assetfs.AssetFS) goview.FileHandler {
	return func(config goview.Config, tplFile string) (content string, err error) {
		// get file contents as string
		filename := tplFile + config.Extension
		data, err := viewsRootBox.Asset(filename)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
}
