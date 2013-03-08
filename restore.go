package main

import (
	"log"
)

//TODO : warning uid change 
//
// TODO : restore channel

func restore_lost_server() {
	sessions := get_restore_sessions()
	log.Print("=== Crash ? Restoring session....")
	//restore server
	for _, session := range sessions {
		_, uid, pseudo, _ := get_user_by_id(session.id_user)
		var keyuser = -1
		if get_key_allusers_by_id(session.id_user) == -1 {
			newid := get_new_id_user()
			us := &User{uid, session.id_user, pseudo, false, make(map[int]*IrcConnec), make(map[int]*Buffer), nil}
			all_users[newid] = us
			keyuser = newid
		} else {
			keyuser = get_key_allusers_by_id(session.id_user)
		}
		go connect_server(session.server, keyuser)
		log.Print("RESTORING : ", pseudo, " reconnecting on ", session.server)
	}

}
