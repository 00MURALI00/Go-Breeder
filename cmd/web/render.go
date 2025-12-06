package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type templateDate struct {
	Data map[string]any
}

func (app *application) render(w http.ResponseWriter, t string, td *templateDate) {
	var tmpl *template.Template

	// if we are using the template chache try to get the template from our map

	if app.config.useCache {
		if templateFromMap, ok := app.templateMap[t]; ok {
			tmpl = templateFromMap
		}
	}

	if tmpl == nil {
		newTemplate, err := app.buildTemplateFromDisk(t)
		if err != nil {
			log.Println("Error building template:", err)
			return
		}
		log.Println("building template foorm disk")
		tmpl = newTemplate
	}

	if td == nil {
		td = &templateDate{}
	}

	if err := tmpl.ExecuteTemplate(w, t, td); err != nil {
		log.Println("Error executing template", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *application) buildTemplateFromDisk(t string) (*template.Template, error) {
	templateSlice := []string{
		"./template/base.layout.gohtml",
		"./template/partials/header.partial.gohtml",
		"./template/partials/footer.partial.gohtml",
		fmt.Sprintf("./template/%s", t),
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		return nil, err
	}

	app.templateMap[t] = tmpl
	return tmpl, nil
}
