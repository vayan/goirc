package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func HandleErrorSql(er error) bool {
	if er != nil {
		log.Println(er)
	}
	return false
}

func test_sql() {
	log.Print("=== TEST SQL ===")
	db, err := sql.Open("mysql", serv_set.DB_user+":"+serv_set.DB_pass+"@("+serv_set.DB_server+":3306)/"+serv_set.DB_name+"?charset=utf8")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Query("SELECT * FROM logirc ")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("=== GOOD ===")
}

func connect_sql() *sql.DB {
	db, err := sql.Open("mysql", serv_set.DB_user+":"+serv_set.DB_pass+"@("+serv_set.DB_server+":3306)/"+serv_set.DB_name+"?charset=utf8")
	if err != nil {
		log.Print(err)
	}
	return db
}

func get_network() {
	db := connect_sql()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM network")
	if err != nil {
		return
	}
	net := make([]*Network, 0, 1)
	var name, adress string
	var id, port, limit int

	for rows.Next() {
		err = rows.Scan(&id, &name, &adress, &port, &limit)
		if err != nil {
			// TODO : Handle error
		}
		net = append(net, &Network{id, name, adress, port, limit, 0})
	}
	network = net
}

// Get backlog from channel
func get_backlog(id_user int, buffer string) []BackLog {
	db := connect_sql()
	defer db.Close()
	//TODO : fix backlog time
	rows, err := db.Query("SELECT nick, message, time FROM logirc WHERE id_user = ? AND buffer = ? ORDER BY time ASC LIMIT 0, 100", id_user, buffer)
	if err != nil {
		return nil
	}
	backlog := make([]BackLog, 0, 1)
	var nick, message, timesql string
	for rows.Next() {
		err = rows.Scan(&nick, &message, &timesql)
		if err != nil {
			// TODO : Handle error
		}
		date, _ := time.Parse("2006-01-02 15:04:05", timesql)
		dateaff := date.Format("15:04")
		backlog = append(backlog, BackLog{nick, message, dateaff})
	}
	return backlog
}

func insert_new_message(id_user int, buffer string, nick string, message string) {
	db := connect_sql()
	defer db.Close()
	_, err := db.Exec("INSERT INTO logirc (id_user, buffer, nick, message, time) VALUES (?, ?, ?, ?, NOW())", id_user, buffer, nick, message)
	HandleErrorSql(err)
}

//get preference for ui
func get_preference() {
	db := connect_sql()
	defer db.Close()
	var name string
	var descr string
	var short_descr string
	var long_descr string
	var base_url string
	var max_lenght_pseudo string

	ar, err := db.Query("SELECT * FROM preference")
	if err != nil {
		return
	}
	ar.Next()
	err = ar.Scan(&name, &descr, &short_descr, &long_descr, &base_url, &max_lenght_pseudo)
	Pref = Preference{name: name, descr: descr, short_descr: short_descr, long_descr: long_descr, base_url: base_url, max_lenght_pseudo: Atoi(max_lenght_pseudo)}
}

//get user settings
//TODO : get settings in sql and update user

//create user setting
func create_settings(user User) {
	db := connect_sql()
	defer db.Close()
	_, err := db.Exec(
		"INSERT INTO users_settings (id_user, notify, save_session) VALUES (?, ?, ?)",
		user.id,
		user.Settings.Notify,
		user.Settings.Save_session)
	HandleErrorSql(err)
}

//update user setting
func update_settings(user User) {
	db := connect_sql()
	defer db.Close()
	var id int

	row := db.QueryRow("SELECT id FROM users_settings WHERE id_user = ? ", user.id)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		create_settings(user)
	} else if err == nil {
		_, err = db.Exec(
			"UPDATE users_settings SET notify = ?, save_session = ? WHERE id_user = ?",
			user.Settings.Notify,
			user.Settings.Save_session,
			user.id)
		HandleErrorSql(err)
	}
}
