package main

import (
	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

const (
	ANON   = 5
	REGIST = 4
)

var (
	all_users = make(map[int]*User)
	db        *sql.DB
	Pref      Preference
	decoder   = schema.NewDecoder()
	store     = sessions.NewCookieStore([]byte("supersecretkeydelamortquitue"))
	serv_set  Server_Settings
)
