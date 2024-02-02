package gosynchro

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed all:static/*
var staticFS embed.FS

//go:embed all:templates/*
var templateFS embed.FS

var StaticFS fs.FS
var ErrorTemplate *template.Template

func init() {
	ErrorTemplate = template.Must(template.ParseFS(templateFS, "templates/error.gohtml"))

	var err error
	StaticFS, err = fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
}
