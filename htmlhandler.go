package main

import (
	"code.google.com/p/go.net/websocket"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

func need_perm(need_defcon int, r *http.Request) bool {
	// defcon 5 = anon
	// defcon 4 = register
	// defcon 1 = admin
	defcon := ANON

	session, _ := store.Get(r, "usersession")
	log.Print("il est ", session.Values["login"])
	if session.Values["login"] == true {
		defcon = REGIST
	}

	if defcon == need_defcon {
		return true
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
	session, _ := store.Get(r, "usersession")

	p := &Page{
		Title: "IRC in your browser",
		Data:  map[string]string{"one": "one"}}

	if session.Values["login"] == true {
		p.Data["login"] = "login"
	}

	RenderHtml(w, "index", p)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage()
	if need_perm(ANON, r) {
		RenderHtml(w, "ajx/register", p)
		return
	}
	HomeHandler(w, r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage()
	if need_perm(ANON, r) {
		RenderHtml(w, "ajx/login", p)
		return
	}
	HomeHandler(w, r)
}

func IrcHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage()
	if need_perm(REGIST, r) {
		RenderHtml(w, "ajx/irc", p)
		return
	}
	HomeHandler(w, r)
}

func ActionLogOut(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "usersession")

	session.Values["login"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func ActionLoginHandler(w http.ResponseWriter, r *http.Request) {

	mail := r.FormValue("InputMail")
	pass := r.FormValue("InputPass")

	session, _ := store.Get(r, "usersession")
	session.Values["login"] = valid_user(mail, pass)
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func ActionRegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	use := new(RegisteringUser)
	decoder.Decode(use, r.Form)
	insert_new_user(*use)
	http.Redirect(w, r, "/", http.StatusFound)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "usersession")
	if session.Values[42] == 43 {
		log.Print("user login")
	}

	p := &Page{
		Title: "IRC in your browser",
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
	r.HandleFunc("/ajx/login", LoginHandler)
	r.HandleFunc("/ajx/irc", IrcHandler)

	//action form
	r.HandleFunc("/register", ActionRegisterHandler)
	r.HandleFunc("/login", ActionLoginHandler)
	r.HandleFunc("/logout", ActionLogOut)

	//websocket
	r.Handle("/ws", websocket.Handler(WsHandle))

	//all js/img/stuff
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))

	log.Println("========Listening on " + port_http + "========")
	log.Fatal(http.ListenAndServe(":"+port_http, r))

}
