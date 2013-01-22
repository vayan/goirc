package main

import (
	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
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
	//TODO : verif pseudo / mail pas deja existant
	if (strings.Contains(user.InputMail, "@")) && (user.InputPass == user.InputPassVerif) && (len(user.InputPseudo) <= Pref.max_lenght_pseudo) {
		_, err := db.Query("INSERT INTO users (pseudo, mail, password) VALUES ('" + user.InputPseudo + "', '" + user.InputMail + "',  '" + EncryptPass(user.InputPass) + "')")
		HandleErrorSql(err)
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
