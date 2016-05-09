package main

import (
    "io/ioutil"
    "net/http"
    "html/template"
    "regexp"
    "errors"
)

//template
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
//invalidation
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

/*
    Simple page which can be save to and load from file
 */
type Page struct {
    Title string
    Body  []byte
}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title:title, Body:body}, nil
}

/*
    Handlers
 */
func viewHandler(w http.ResponseWriter, r *http.Request) {
    //title := r.URL.Path[len("/view/"):]
    //validate get title
    title, err:= getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/" + title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    //title := r.URL.Path[len("/edit/"):]
    //validate get title
    title, err:= getTitle(w, r)
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
    //title := r.URL.Path[len("/save/"):] //validate get title
    title, err:= getTitle(w, r)
    if err != nil {
        return
    }
    body := r.FormValue("body")
    p := &Page{Title:title, Body:[]byte(body)}
    err = p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    //origin template
    //t, err := template.ParseFiles(tmpl + ".html")
    //if err != nil {
    //    http.Error(w, err.Error(), http.StatusInternalServerError)
    //    return
    //}
    //err = t.Execute(w, p)
    //if err != nil {
    //    http.Error(w, err.Error(), http.StatusInternalServerError)
    //}

    //template with cache
    err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    //test for simple save to and load from file
    //p1 := &Page{Title:"TestPage", Body:[]byte("This is a sample Page.")}
    //p1.save()
    //p2, _ := loadPage("TestPage")
    //fmt.Println(string(p2.Body))

    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    http.ListenAndServe(":8101", nil)
}


/*
 simple http server
 */
//func handler(w http.ResponseWriter, r *http.Request) {
//    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
//}
//
//func main() {
//    http.HandleFunc("/", handler)
//    http.ListenAndServe(":8001", nil)
//}
