
package main

import (
	"io/ioutil"
	"path"
	"os"
	"log"
	"net/http"
	"html/template"
	"fmt"
  "regexp"
  "errors"
)

/***************
 * helpers
 **************/

func getPageFileNameWithPath(title string) string {
	return path.Join(pagesDir, title+".txt")
}

func getTemplateFileNameWithoutPath(tmplName string) string {
  return tmplName + ".html"
}

func getTemplateFileNameWithPath(tmplName string) string {
	return path.Join(templatesDir, getTemplateFileNameWithoutPath(tmplName))
}

func createDirs() {
	dirs := [...]string{pagesDir, templatesDir}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func cacheTemplates() *template.Template {
  v := getTemplateFileNameWithPath("view")
  e := getTemplateFileNameWithPath("edit")
  d := getTemplateFileNameWithPath("default")
  return template.Must(template.ParseFiles(v, e, d))
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    msg := "invalid url: "+ r.URL.Path
    fmt.Println(msg)
    return "", errors.New(msg)
  }
  return m[2], nil  //<- m: 0=url, 1=handler, 2=title
}

func renderTemplate(w http.ResponseWriter, tmplName string, p *Page) {
  err := templates.ExecuteTemplate(w, getTemplateFileNameWithoutPath(tmplName), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/***************
 * Page
 **************/

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	return ioutil.WriteFile(getPageFileNameWithPath(p.Title), p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	body, err := ioutil.ReadFile(getPageFileNameWithPath(title))
	if err != nil {
		return nil, err
	}
	return &Page{Title:title, Body:body}, nil
}

/***************
 * handlers
 **************/

func viewHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("in viewHandler. url: %s\n", r.URL.Path)
	title, err := getTitle(w, r)
  if err != nil {
    return
  }
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("in editHandler. url: %s\n", r.URL.Path)
	title, err := getTitle(w, r)
  if err != nil {
    return
  }
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title:title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("in saveHandler. url: %s\n", r.URL.Path)
  title, err := getTitle(w, r)
  if err != nil {
    return
  }
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("in deleteHandler. url: %s\n", r.URL.Path)
  title, err := getTitle(w, r)
  if err != nil {
    return
  }
	fileName := getPageFileNameWithPath(title)
	if _, err := os.Stat(fileName); err == nil {
		err = os.Remove(fileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		fmt.Printf("doesn't exist: %s\n", fileName)
	}
  http.Redirect(w, r, "/", http.StatusFound)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("in defaultHandler. url: %s\n", r.URL.Path)
  p := &Page{Title:"wiki", Body:[]byte(r.URL.Path)}
	renderTemplate(w, "default", p)
}

/***************
 * main
 **************/

var pagesDir     = "pages"
var templatesDir = "templates"
var templates    = cacheTemplates()
var routeFilter  = "^(/|/view/|/edit/|/save/|/delete/)([a-zA-Z0-9_]*)$"
var validPath    = regexp.MustCompile(routeFilter)

func main() {
	createDirs()

	http.HandleFunc("/view/",   viewHandler)
	http.HandleFunc("/edit/",   editHandler)
	http.HandleFunc("/save/",   saveHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/",        defaultHandler)

  addr := ":8080"
  fmt.Printf("listening on %s\n", addr)
  err := http.ListenAndServe(addr, nil)
  if err != nil {
    log.Fatal(err)
  }
}













