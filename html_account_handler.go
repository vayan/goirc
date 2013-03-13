package main

import (
	_ "github.com/Go-SQL-Driver/MySQL"
	"log"
	"net/http"
)

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

func ActionLogOut(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, serv_set.Cookie_session)

	session.Values["login"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func ActionLoginHandler(w http.ResponseWriter, r *http.Request) {

	mail := r.FormValue("InputMail")
	pass := r.FormValue("InputPass")

	session, _ := store.Get(r, serv_set.Cookie_session)
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
