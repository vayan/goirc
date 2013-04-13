package main

import (
	"code.google.com/p/go.net/websocket"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

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
	r.HandleFunc("/ajx/settings", SettingsHandler)
	r.HandleFunc("/ajx/getsettings", GetSettingsHandler)
	r.HandleFunc("/ajx/getfriends", GetFriendsHandler)
	r.HandleFunc("/ajx/profile", ProfileHandler)
	r.HandleFunc("/ajx/set-channels", SetChanHandler)

	//action form
	r.HandleFunc("/register", ActionRegisterHandler)
	r.HandleFunc("/login", ActionLoginHandler)
	r.HandleFunc("/logout", ActionLogOut)

	//websocket
	r.Handle("/ws", websocket.Handler(WsHandle))

	//all js/img/stuff
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))

	log.Print("=== Listening on " + serv_set.Port_http + " ===")
	log.Fatal(http.ListenAndServe(":"+serv_set.Port_http, r))
}
