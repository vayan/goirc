package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"database/sql"
	"log"
)

func HandleErrorSql(er error) bool {
	if er != nil {
		log.Println(er)
	}
	return false
}

func connect_sql() {
	log.Println("========Start connexion DB========")
	var err error

	db, err = sql.Open("mysql", DB_USER+":"+DB_PASS+"@("+DB_SERV+":3306)/"+DB_NAME+"?charset=utf8")
	HandleErrorSql(err)
}

func insert_new_user(user RegisteringUser) {
	//ar, err := db.Query("SELECT * FROM user WHERE pseudo='" + user.InputPseudo + "' OR email='" + user.InputMail + "' ")
	//HandleErrorSql(err)
}

func get_preference() {
	var name string
	var descr string
	var short_descr string
	var long_descr string
	var base_url string

	ar, err := db.Query("SELECT * FROM preference")
	HandleErrorSql(err)
	ar.Next()
	err = ar.Scan(&name, &descr, &short_descr, &long_descr, &base_url)
	Pref = Preference{name: name, descr: descr, short_descr: short_descr, long_descr: long_descr, base_url: base_url}
}
