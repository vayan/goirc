package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

const (
	root_web = "templates/"
)

type Page struct {
	Title string
	Data  map[string]string
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
	p := &Page{Title: "Home", Data: get_preference()}
	RenderHtml(w, "ajx/home", p)
}

func start_http_server() {
	fmt.Print("\n\n========Start goric web server========\n")
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/ajx/home", HomeHandler)
	r.HandleFunc("/ajx/register", RegisterHandler)
	r.HandleFunc("/ajx/irc", IrcHandler)
	r.Handle("/ws", websocket.Handler(WsHandle))
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))
	http.ListenAndServe(":1111", r)
}
