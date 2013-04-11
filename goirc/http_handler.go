package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func need_perm(need_defcon int, r *http.Request) bool {
	defcon := ANON

	session, _ := store.Get(r, serv_set.Cookie_session)
	if session.Values["login"] == true {
		defcon = REGIST
	}

	if defcon == need_defcon {
		return true
	}
	return false
}

func loadPage(title string) *Page {
	return &Page{Title: title}
}

func RenderHtml(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(serv_set.Root_web + tmpl + ".html")
	t.Execute(w, p)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, serv_set.Cookie_session)

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

func ActionRegisterHandler(w http.ResponseWriter, r *http.Request) {
	//TODO : captcha
	r.ParseForm()
	use := new(RegisteringUser)
	decoder.Decode(use, r.Form)
	_, ret := insert_new_user(*use)
	ret_json := make(map[string]([]string))
	ret_json["errors"] = ret
	b, _ := json.Marshal(ret_json)
	fmt.Fprint(w, string(b))
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

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	if need_perm(REGIST, r) {
		if r.Method == "POST" {
			session, _ := store.Get(r, serv_set.Cookie_session)
			user := get_user_id(session.Values["id"].(int))
			if notify, err := strconv.ParseBool(r.FormValue("Notify")); err == nil {
				user.Settings.Notify = notify
			}
			if save_session, err := strconv.ParseBool(r.FormValue("Save_Session")); err == nil {
				user.Settings.Save_session = save_session
			}
			update_settings(*user)
		}
		p := loadPage("Settings")
		RenderHtml(w, "ajx/settings", p)
		return
	}
	HomeHandler(w, r)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if need_perm(REGIST, r) {
		session, _ := store.Get(r, serv_set.Cookie_session)
		//user := get_user_id(session.Values["id"].(int))
		p := &Page{
			Title: "Profile",
			Data: map[string]interface{}{
				"hashmail": get_mail_hash(session.Values["mail"].(string)),
				"pseudo":   session.Values["pseudo"].(string)}}
		RenderHtml(w, "ajx/profile", p)
		return
	}
	HomeHandler(w, r)
}
