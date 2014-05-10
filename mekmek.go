package main 

import (
	"io/ioutil"
	"net/http"
	"html/template"
)


type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := "db/" + p.Title
	return ioutil.WriteFile(filename,p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := "db/"+ title
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w,p)
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
	  http.Redirect(w,r,"/edit/" + title, http.StatusFound)
	  return
	}
	renderTemplate(w, "view", p)
}

func listHandler(w http.ResponseWriter, r *http.Request){
  listDB, _ := ioutil.ReadDir("db")
  t, _ := template.ParseFiles("list.html")
  t.Execute(w, map[string]interface{} {"ListFile" : listDB})
}

func editHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{ Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request){
  title := r.URL.Path[len("/save/"):]
  body := r.FormValue("body")
  p := &Page{ Title: title, Body: []byte(body)}
  p.save()
  http.Redirect(w,r,"/view/" + title, http.StatusFound)
}


func main() {
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
