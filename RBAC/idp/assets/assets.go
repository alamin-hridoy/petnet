package assets

import (
	"embed"
	"io/fs"
)

//go:embed css/app.min.css fonts images js
var Static embed.FS

//go:embed templates/*.html
var tmpl embed.FS

var Templates fs.FS

func init() {
	t, err := fs.Sub(tmpl, "templates")
	if err != nil {
		Templates = tmpl
	}
	Templates = t
}
