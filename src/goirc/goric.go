package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

const (
	DB_NAME  = "goirc"
	DB_USER  = "goirc"
	DB_PASS  = "rRfCKB6eMnDXNVZw"
	DB_SERV  = "69.85.88.161"
	root_web = "templates/"
)

type Page struct {
	Title string
}

func HandleErrorFatal(er error) bool {
	if er != nil {
		log.Fatal(er)
	}
	return false
}

func loadPage() *Page {
	title := "test"
	return &Page{Title: title}
}

func RenderHtml(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(root_web + tmpl + ".html")
	t.Execute(w, p)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage()
	RenderHtml(w, "index", p)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage()
	RenderHtml(w, "ajx/register", p)
}

func IrcHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage()
	RenderHtml(w, "ajx/irc", p)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage()
	RenderHtml(w, "ajx/home", p)
}

func main() {
	fmt.Print("========Start connexion DB========\n\n")
	//username:password@hostspec/database
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASS+"@("+DB_SERV+":3306)/"+DB_NAME+"?charset=utf8")
	HandleErrorFatal(err)
	ar, err := db.Query("SELECT * FROM preference")
	HandleErrorFatal(err)
	fmt.Print(ar.Columns())

	fmt.Print("\n\n========Start goric web server========\n")
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/ajx/home", HomeHandler)
	r.HandleFunc("/ajx/register", RegisterHandler)
	r.HandleFunc("/ajx/irc", IrcHandler)
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))
	http.ListenAndServe(":1111", r)
}
