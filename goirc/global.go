package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

const (
	ANON   = 5
	REGIST = 4
)

var (
	all_users = make(map[int]*User)
	Pref      Preference
	decoder   = schema.NewDecoder()
	store     = sessions.NewCookieStore([]byte("supersecretkeydelamortquitue"))
	serv_set  Server_Settings
)
