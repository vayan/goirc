package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

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

func ActionRegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	use := new(RegisteringUser)
	decoder.Decode(use, r.Form)
	insert_new_user(*use)
	http.Redirect(w, r, "/", http.StatusFound)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{
		Title: "Home",
		Data:  map[string]string{"name": Pref.name, "descr": Pref.descr, "short_descr": Pref.short_descr, "long_descr": Pref.long_descr, "base_url": Pref.base_url}}
	RenderHtml(w, "ajx/home", p)
}

func start_http_server() {
	log.Println("========Starting goric web server========")
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)

	//ajx html
	r.HandleFunc("/ajx/home", HomeHandler)
	r.HandleFunc("/ajx/register", RegisterHandler)
	r.HandleFunc("/ajx/irc", IrcHandler)

	//action form
	r.HandleFunc("/register", ActionRegisterHandler)

	//wesocket
	r.Handle("/ws", websocket.Handler(WsHandle))

	//all js/img/stuff
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))

	log.Println("========Listening on " + port_http + "========")
	log.Fatal(http.ListenAndServe(":"+port_http, r))

}
