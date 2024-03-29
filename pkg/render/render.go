package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/shiibs/bookings/pkg/config"
	"github.com/shiibs/bookings/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets  the config for template package
func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(tmplData *models.TemplateData) *models.TemplateData {

	return tmplData
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, tmpl string, data *models.TemplateData) {
	var tmplCache map[string]*template.Template

	// get the template cache from the app config
	if app.UseCache {
		tmplCache = app.TemplateCache
	} else {
		tmplCache, _ = CreateTemplateCache()
	}

	// get the requested template from cache
	t, ok := tmplCache[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	data = AddDefaultData(data)

	// render the template
	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("error parsing template", err)
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)

		tmplSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			tmplSet, err = tmplSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = tmplSet

	}
	return myCache, nil
}
