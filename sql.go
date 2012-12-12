package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"database/sql"
	"log"
	"strconv"
	"strings"
)

func Atoi(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
	}
	return val
}

func HandleErrorSql(er error) bool {
	if er != nil {
		log.Println(er)
	}
	return false
}

func connect_sql() {
	log.Println("========Start Connexion DB========")
	var err error
	db, err = sql.Open("mysql", DB_USER+":"+DB_PASS+"@("+DB_SERV+":3306)/"+DB_NAME+"?charset=utf8")
	HandleErrorSql(err)
}

func insert_new_user(user RegisteringUser) int {
	//a finir insert user db
	if !(strings.Contains(user.InputMail, "@")) {
		log.Println("Erreur mail")
	}
	if !(user.InputPass == user.InputPassVerif) {
		log.Println("Erreur pass")
	}
	if !(len(user.InputPseudo) <= Pref.max_lenght_pseudo) {
		log.Println("Erreur pseudo")
	}
	return -1
}

func get_preference() {
	var name string
	var descr string
	var short_descr string
	var long_descr string
	var base_url string
	var max_lenght_pseudo string

	ar, err := db.Query("SELECT * FROM preference")
	HandleErrorSql(err)
	ar.Next()
	err = ar.Scan(&name, &descr, &short_descr, &long_descr, &base_url, &max_lenght_pseudo)
	Pref = Preference{name: name, descr: descr, short_descr: short_descr, long_descr: long_descr, base_url: base_url, max_lenght_pseudo: Atoi(max_lenght_pseudo)}
}
