package main

import (
	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

var all_users = make(map[int]*User)
var db *sql.DB
var Pref Preference
var decoder = schema.NewDecoder()
var store = sessions.NewCookieStore([]byte("supersecretkeydelamortquitue"))
var set Settings
