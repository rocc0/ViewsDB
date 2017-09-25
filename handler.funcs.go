package main

import (
	"html/template"
)

var TemplateHelpers = template.FuncMap{
	"toString": func(s []uint8) string {
		return string(s)
	},
}
