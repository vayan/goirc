package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"database/sql"
	"github.com/gorilla/schema"
)

var all_users = make(map[int]*User)
var db *sql.DB
var Pref Preference
var decoder = schema.NewDecoder()
