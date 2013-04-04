package main

import (
	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
	"html/template"
	"log"
	"time"
)

//TODO : Check loose connexion
// TODO : timeout sql deco

func HandleErrorSql(er error) bool {
	if er != nil {
		log.Println(er)
	}
	return false
}

func connect_sql() {
	log.Println("=== Start Connexion DB ===")
	var err error
	db, err = sql.Open("mysql", serv_set.DB_user+":"+serv_set.DB_pass+"@("+serv_set.DB_server+":3306)/"+serv_set.DB_name+"?charset=utf8")
	HandleErrorSql(err)
}

// Get backlog from channel
func get_backlog(id_user int, buffer string) []*BackLog {
	rows, err := db.Query("SELECT nick, message, time FROM logirc WHERE id_user = ? AND buffer = ? ORDER BY time ASC", id_user, buffer)
	HandleErrorSql(err)
	backlog := make([]*BackLog, 0, 10)
	var nick, message, timesql string
	for rows.Next() {
		err = rows.Scan(&nick, &message, &timesql)
		if err != nil {
			// TODO : Handle error
		}
		date, _ := time.Parse("2006-01-02 15:04:05", timesql)
		dateaff := date.Format("15:04")
		backlog = append(backlog, &BackLog{nick, message, dateaff})
	}
	return backlog
}

func insert_new_message(id_user int, buffer string, nick string, message string) {
	message = template.HTMLEscapeString(message)
	_, err := db.Exec("INSERT INTO logirc (id_user, buffer, nick, message, time) VALUES (?, ?, ?, ?, NOW())", id_user, buffer, nick, message)
	HandleErrorSql(err)
}

//get preference for ui
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

//get user settings
//TODO : get settings in sql and update user

//create user setting
func create_settings(user User) {
	_, err := db.Exec(
		"INSERT INTO users_settings (id_user, notify, save_session) VALUES (?, ?, ?)",
		user.id,
		user.Settings.Notify,
		user.Settings.Save_session)
	HandleErrorSql(err)
}

//update user setting
func update_settings(user User) {
	var id int

	row := db.QueryRow("SELECT id FROM users_settings WHERE id_user = ? ", user.id)
	err := row.Scan(&id)
	HandleErrorSql(err)
	if err == sql.ErrNoRows {
		create_settings(user)
	} else {
		_, err = db.Exec(
			"UPDATE users_settings SET notify = ?, save_session = ? WHERE id_user = ?",
			user.Settings.Notify,
			user.Settings.Save_session,
			user.id)
		HandleErrorSql(err)
	}
}
