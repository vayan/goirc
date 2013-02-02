package main

import (
	"code.google.com/p/go.net/websocket"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

// FIXME : random deco..

func need_perm(need_defcon int, r *http.Request) bool {
	// defcon 5 = anon
	// defcon 4 = register
	// defcon 1 = admin
	defcon := ANON

	session, _ := store.Get(r, COOKIE_SESSION)
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
	session, _ := store.Get(r, COOKIE_SESSION)

	if session.Values["login"] == nil {
		session.Values["login"] = false
	}

	session.Values["uid"] = ""
	if session.Values["login"].(bool) {
		log.Print("User connected get UID from bdd")
		session.Values["uid"] = get_uid(session.Values["id"].(int))
	}
	session.Save(r, w)
	uid := session.Values["uid"]

	p := &Page{
		Title: "IRC in your browser",
		Data: map[string]string{
			"one": "one",
			"uid": uid.(string)}}

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
	session, _ := store.Get(r, COOKIE_SESSION)

	session.Values["login"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func ActionLoginHandler(w http.ResponseWriter, r *http.Request) {

	mail := r.FormValue("InputMail")
	pass := r.FormValue("InputPass")

	session, _ := store.Get(r, COOKIE_SESSION)
	valid, id, pseudo, email, uid := get_user(mail, pass)
	session.Values["login"] = valid
	session.Values["id"] = id
	session.Values["pseudo"] = pseudo
	session.Values["mail"] = email
	if len(uid) < 1 {
		log.Print("UID never generate, generate new uid")
		uid = generate_unique_uid(pseudo)
		set_uid(id, uid)
	}
	session.Values["uid"] = uid
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
	//session, _ := store.Get(r, COOKIE_SESSION)

	p := &Page{
		Title: "IRC in your browser",
		Data: map[string]string{
			"name":        Pref.name,
			"descr":       Pref.descr,
			"short_descr": Pref.short_descr,
			"long_descr":  Pref.long_descr,
			"base_url":    Pref.base_url}}
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
