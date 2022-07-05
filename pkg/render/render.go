package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/sureshkumawat/bookings/pkg/config"
	"github.com/sureshkumawat/bookings/pkg/models"
)

var functions = template.FuncMap{}
var basic_path = "/Users/Suresh Kumawat/Desktop/Coding parts/go_coding/go_course/go_hello/"
var app *config.AppConfig

//NewTemplates sets the config to new template package
func NewTemplates(a *config.AppConfig) {
	app = a

}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreatTemplateCache()
	}
	t, ok := tc[tmpl]
	fmt.Println(t)
	if !ok {
		log.Fatal("could not get template from template cache")
	}
	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}
}

//CreatTemplateCache creates a template cache as a map
func CreatTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(basic_path + "templates/*.page.html")
	//fmt.Println(pages)
	if err != nil {

		return myCache, err
	}
	//fmt.Println("inside test functon")
	for _, page := range pages {
		name := filepath.Base(page)
		//fmt.Println(name)
		//fmt.Println("page is currently", name)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println("there is an error2")
			return myCache, err
		}

		matches, err := filepath.Glob(basic_path + "templates/*.layout.html")
		if err != nil {
			fmt.Println("there is an error3")
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(basic_path + "templates/*.layout.html")
			if err != nil {
				fmt.Println("there is an error4")
				return myCache, err
			}

		}
		myCache[name] = ts
	}
	return myCache, nil
}
