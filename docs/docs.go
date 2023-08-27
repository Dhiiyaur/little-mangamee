package docs

import (
	"embed"
	"html/template"
)

//go:embed index.html
var HtmlBase embed.FS

//go:embed mangamee_collection.yml
var Colection embed.FS

func GetTemplates() (*template.Template, error) {
	allTemplates, err := template.ParseFS(HtmlBase, "docs/*.html")
	if err != nil {
		return nil, err
	}
	return allTemplates, nil
}
