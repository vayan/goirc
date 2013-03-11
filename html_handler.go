package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// FIXME : random deco..

func need_perm(need_defcon int, r *http.Request) bool {
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
		Data: map[string]interface{}{
			"one": "one",
			"uid": uid.(string)}}

	if session.Values["login"] == true {
		p.Data["login"] = "login"
	}

	RenderHtml(w, "index", p)
}

func IrcHandler(w http.ResponseWriter, r *http.Request) {
	HomeHandler(w, r)
}

func SetServHandler(w http.ResponseWriter, r *http.Request) {
	//TODO : check connected
	if need_perm(REGIST, r) {
		p := loadPage()
		RenderHtml(w, "ajx/set-servers", p)
	}
}

func SetChanHandler(w http.ResponseWriter, r *http.Request) {
	if need_perm(REGIST, r) {
		//TODO : test si ws active

		var allserv string
		session, _ := store.Get(r, COOKIE_SESSION)
		us := get_user_id(session.Values["id"].(int))

		for _, irco := range us.Buffers {
			if irco.name[0] != '#' {
				allserv += "<option value='" + strconv.Itoa(irco.id) + "'>" + irco.name + "</option>"
			}
		}
		p := &Page{
			Title: "IRC in your browser",
			Data:  map[string]interface{}{"servers": template.HTML(allserv)}}
		RenderHtml(w, "ajx/set-channels", p)
	}
}

func ActionBacklogHandler(w http.ResponseWriter, r *http.Request) {
	//TODO : JSON all that
	idbuffer := Atoi(r.FormValue("idbuffer"))
	session, _ := store.Get(r, COOKIE_SESSION)

	if need_perm(REGIST, r) {
		user := get_user_id(session.Values["id"].(int))
		buffers := user.Buffers
		for _, buff := range buffers {
			if buff.id == idbuffer {
				backlog := get_backlog(user.id, user.Buffers[idbuffer].addr)
				for _, log := range backlog {
					fmt.Fprint(w, "<tr class=\"msg\"><td class=\"pseudo "+log.nick+"\">"+log.nick+"</td><td class=\"message\"><div class='messagediv'>"+log.message+"</div></td><td class=\"time\">"+log.time+"</td></tr>")
				}
				return
			}
		}

	}
}

func ActionRegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	use := new(RegisteringUser)
	decoder.Decode(use, r.Form)
	insert_new_user(*use)
	http.Redirect(w, r, "/", http.StatusFound)
}

func UsersListHandler(w http.ResponseWriter, r *http.Request) {
	//TODO : check connected

	jsonres := "{ \"UserList\":["
	id := Atoi(r.FormValue("channel"))
	session, _ := store.Get(r, COOKIE_SESSION)
	us := get_user_id(session.Values["id"].(int))

	if _, ok := us.Buffers[id]; ok {
		for e := us.Buffers[id].users.Front(); e != nil; e = e.Next() {
			b, _ := json.Marshal(e.Value.(ChannelUser))
			jsonres += string(b) + ","
		}
		jsonres = jsonres[:len(jsonres)-1]
		jsonres += "]}"
		fmt.Fprint(w, jsonres)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{
		Title: "IRC in your browser",
		Data: map[string]interface{}{
			"name":        Pref.name,
			"descr":       Pref.descr,
			"short_descr": Pref.short_descr,
			"long_descr":  Pref.long_descr,
			"base_url":    Pref.base_url}}
	RenderHtml(w, "ajx/home", p)
}

func start_http_server() {
	log.Print("=== Starting goric web server ===")
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)

	//ajx html
	r.HandleFunc("/ajx/backlog", ActionBacklogHandler)
	r.HandleFunc("/ajx/home", HomeHandler)
	r.HandleFunc("/ajx/register", RegisterHandler)
	r.HandleFunc("/ajx/login", LoginHandler)
	r.HandleFunc("/ajx/irc", IrcHandler)
	r.HandleFunc("/ajx/userslist", UsersListHandler)

	//ajx html settings
	r.HandleFunc("/ajx/set-servers", SetServHandler)
	r.HandleFunc("/ajx/set-channels", SetChanHandler)

	//action form
	r.HandleFunc("/register", ActionRegisterHandler)
	r.HandleFunc("/login", ActionLoginHandler)
	r.HandleFunc("/logout", ActionLogOut)

	//websocket
	r.Handle("/ws", websocket.Handler(WsHandle))

	//all js/img/stuff
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))

	log.Print("=== Listening on " + port_http + " ===")
	log.Fatal(http.ListenAndServe(":"+port_http, r))
}
