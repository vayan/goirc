package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"log"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage("Register")
	if need_perm(ANON, r) {
		RenderHtml(w, "ajx/register", p)
		return
	}
	HomeHandler(w, r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage("Login")
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
	//TODO : limit number of try and captcha
	ret := make(map[string]([]string))
	ret["errors"] = make([]string, 10)
	session, _ := store.Get(r, serv_set.Cookie_session)
	mail := r.FormValue("InputMail")
	pass := r.FormValue("InputPass")

	if len(mail) < 1 || len(pass) < 1 {
		ret["errors"][0] = "Please fill all fields"
	}
	valid, id, pseudo, email, uid := get_user(mail, pass)
	if !valid {
		ret["errors"][1] = "Incorrect pseudo or password"
	}
	session.Values["login"] = valid
	session.Values["id"] = id
	session.Values["pseudo"] = pseudo
	session.Values["mail"] = email
	if len(uid) < 1 {
		log.Print("UID never generate, generate new uid")
		uid = generate_unique_uid(email)
		set_uid(id, uid)
	}
	session.Values["uid"] = uid
	session.Save(r, w)
	b, _ := json.Marshal(ret)
	fmt.Fprint(w, string(b))
}
