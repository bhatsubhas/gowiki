package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := getFilename(p.Title)
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := getFilename(title)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := getTitle(r.URL.Path, "/view/")
	log.Printf("Request received to view title - %s", title)
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "templates/view.html", page)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := getTitle(r.URL.Path, "/edit/")
	log.Printf("Request received to edit title - %s", title)
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	renderTemplate(w, "templates/edit.html", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := getTitle(r.URL.Path, "/save/")
	log.Printf("Request received to save title - %s", title)
	body := r.FormValue("body")
	page := &Page{Title: title, Body: []byte(body)}
	page.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, templ string, p *Page) {
	t, _ := template.ParseFiles(templ)
	t.Execute(w, p)
}

func getTitle(path string, offset string) string {
	return path[len(offset):]
}

func getFilename(title string) string {
	return title + ".txt"
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Printf("Starting the web server in port %s", "8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
