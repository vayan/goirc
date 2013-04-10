package main

import (
	_ "github.com/Go-SQL-Driver/MySQL"
	"log"
	"strings"
)

// Get user UID by id
func get_uid(id int) string {
	var uid string

	row := db.QueryRow("SELECT uid FROM users WHERE id = ?", id)
	err := row.Scan(&uid)
	if err != nil {
		log.Println(err)
	}
	return uid
}

//Set UID user
func set_uid(id int, uid string) {
	_, err := db.Exec("UPDATE users SET uid = ? WHERE id = ?", uid, id)
	HandleErrorSql(err)
}

//return user by uid
// TODO : return struct user
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

//TODO : return struct
func get_user_by_id(id int) (bool, string, string, string) {
	var pseudo, mail, uid string

	valid := false
	row := db.QueryRow("SELECT uid, pseudo, mail FROM users WHERE id = ? ", id)

	err := row.Scan(&uid, &pseudo, &mail)
	if len(pseudo) > 0 {
		valid = true
	}
	if err != nil {
		log.Println(err)
	}
	return valid, uid, pseudo, mail
}

//get user from mail / pass
func get_user(email string, pass string) (bool, int, string, string, string) {
	var pseudo, mail, uid string
	var id int

	row := db.QueryRow("SELECT id, pseudo, mail, uid FROM users WHERE mail = ? AND password = ? ", email, EncryptPass(pass))
	err := row.Scan(&id, &pseudo, &mail, &uid)
	if err != nil {
		return false, id, pseudo, mail, uid
	}
	return true, id, pseudo, mail, uid
}

func user_exist(mail string) bool {
	row := db.QueryRow("SELECT mail FROM users WHERE mail = ?", mail)
	err := row.Scan(&mail)
	if err != nil {
		return false
	}
	return true
}

//get session to restore (session with channel)
func get_restore_sessions() []*RestoreSession {
	rows, err := db.Query("SELECT id_user, server, channel, friends FROM session_save")
	HandleErrorSql(err)
	restoresessions := make([]*RestoreSession, 0, 10)
	var server, channel, friends string
	var id_user int
	for rows.Next() {
		err = rows.Scan(&id_user, &server, &channel, &friends)
		if err != nil {
			// TODO : Handle error
		}
		restoresessions = append(restoresessions, &RestoreSession{id_user, server, channel, friends})
	}
	return restoresessions
}

func insert_new_friend_session(id_user int, server string, friends string) {
	var ffriends string

	row := db.QueryRow("SELECT friends FROM session_save WHERE id_user = ? AND server = ? ", id_user, server)
	err := row.Scan(&ffriends)
	for _, val := range strings.Split(ffriends, ",") {
		if val == friends {
			return
		}
	}
	friends = "," + friends
	_, err = db.Exec("UPDATE session_save SET friends = CONCAT(friends, ?) WHERE id_user = ? AND server = ?", friends, id_user, server)
	HandleErrorSql(err)
}

func insert_new_server_session(id_user int, server string) {
	var sserver string

	row := db.QueryRow("SELECT server FROM session_save WHERE id_user = ? AND server = ? ", id_user, server)
	err := row.Scan(&sserver)
	if len(sserver) < 1 {
		_, err = db.Exec("INSERT INTO session_save (id_user, server, channel) VALUES (?, ?, ?)", id_user, server, "")
		HandleErrorSql(err)
	}
}

func insert_new_channel_session(id_user int, server string, channel string) {
	var cchanel string

	row := db.QueryRow("SELECT channel FROM session_save WHERE id_user = ? AND server = ? ", id_user, server)
	err := row.Scan(&cchanel)
	for _, val := range strings.Split(cchanel, ",") {
		if val == channel {
			return
		}
	}
	channel = "," + channel
	_, err = db.Exec("UPDATE session_save SET channel = CONCAT(channel, ?) WHERE id_user = ? AND server = ?", channel, id_user, server)
	HandleErrorSql(err)
}

func insert_new_user(user RegisteringUser) (int, []string) {
	msg_err := make([]string, 10)
	ret := 0
	// TODO : welcom mail to send
	if user_exist(user.InputMail) {
		msg_err[0] = "Email already used"
		ret = 1
	}
	if !strings.Contains(user.InputMail, "@") {
		msg_err[1] = "Incorrect Email"
		ret = 1
	}
	if user.InputPass != user.InputPassVerif {
		msg_err[2] = "Passwords don't match"
		ret = 1
	}
	//TODO : min lenght not max
	if len(user.InputPseudo) > Pref.max_lenght_pseudo {
		msg_err[3] = "Passwords too short"
		ret = 1
	}
	if ret == 0 {
		cleanpseudo := strings.Trim(strings.ToLower(user.InputPseudo), " ")
		cleanmail := strings.Trim(strings.ToLower(user.InputMail), " ")
		cleanpass := strings.Trim(user.InputPass, " ")
		_, err := db.Exec("INSERT INTO users (pseudo, mail, password) VALUES (?, ?, ?)", cleanpseudo, cleanmail, EncryptPass(cleanpass))
		HandleErrorSql(err)
	}
	return ret, msg_err
}

func get_stats_user(id_user int) UserStats {
	//TODO : maybe..
	return UserStats{1, 1, "", 1}
}
