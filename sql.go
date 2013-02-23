package main

import (
	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
	"log"
	"strings"
)

//TODO : Check loose connexion

func HandleErrorSql(er error) bool {
	if er != nil {
		log.Println(er)
	}
	return false
}

func get_uid(id int) string {
	var uid string

	row := db.QueryRow("SELECT uid FROM users WHERE id = ?", id)
	err := row.Scan(&uid)
	if err != nil {
		log.Println(err)
	}
	return uid
}

func set_uid(id int, uid string) {
	_, err := db.Exec("UPDATE users SET uid = ? WHERE id = ?", uid, id)
	HandleErrorSql(err)
}

func connect_sql() {
	log.Println("=== Start Connexion DB ===")
	var err error
	db, err = sql.Open("mysql", DB_USER+":"+DB_PASS+"@("+DB_SERV+":3306)/"+DB_NAME+"?charset=utf8")
	HandleErrorSql(err)
}

func get_user_by_uid(uid string) (bool, int, string, string) {
	var pseudo, mail string
	var id int

	valid := false
	row := db.QueryRow("SELECT id, pseudo, mail FROM users WHERE uid = ? ", uid)

	err := row.Scan(&id, &pseudo, &mail)
	if len(pseudo) > 0 {
		valid = true
	}
	if err != nil {
		log.Println(err)
	}
	return valid, id, pseudo, mail
}

func get_backlog(id_user int, buffer string) []*BackLog {
	rows, err := db.Query("SELECT nick, message, time FROM logirc WHERE id_user = ? AND buffer = ? ORDER BY time ASC", id_user, buffer)
	HandleErrorSql(err)
	backlog := make([]*BackLog, 0, 10)
	var nick, message, time string
	for rows.Next() {
		err = rows.Scan(&nick, &message, &time)
		if err != nil {
			// TODO : Handle error
		}
		backlog = append(backlog, &BackLog{nick, message, time})
	}
	return backlog
}

func get_user(email string, pass string) (bool, int, string, string, string) {
	var pseudo, mail, uid string
	var id int

	valid := false
	row := db.QueryRow("SELECT id, pseudo, mail, uid FROM users WHERE mail = ? AND password = ? ", email, EncryptPass(pass))

	err := row.Scan(&id, &pseudo, &mail, &uid)
	if len(pseudo) > 0 {
		valid = true
	}
	if err != nil {
		log.Println(err)
	}
	return valid, id, pseudo, mail, uid
}

func insert_new_user(user RegisteringUser) int {
	//TODO : verif pseudo / mail pas deja existant
	// TODO : welcom mail to send

	if (strings.Contains(user.InputMail, "@")) && (user.InputPass == user.InputPassVerif) && (len(user.InputPseudo) <= Pref.max_lenght_pseudo) {
		_, err := db.Exec("INSERT INTO users (pseudo, mail, password) VALUES (?, ?, ?)", user.InputPseudo, user.InputMail, EncryptPass(user.InputPass))
		HandleErrorSql(err)
	}
	return -1
}

func insert_new_message(id_user int, buffer string, nick string, message string) {
	_, err := db.Exec("INSERT INTO logirc (id_user, buffer, nick, message, time) VALUES (?, ?, ?, ?, NOW())", id_user, buffer, nick, message)
	HandleErrorSql(err)
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
