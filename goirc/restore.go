package main

import (
	"log"
	"strings"
)

//TODO : warning uid change

func restore_lost_server() {
	if serv_set.Restore_session == false {
		return
	}
	sessions := get_restore_sessions()
	log.Print("=== Crash ? Restoring session....")
	for _, session := range sessions {
		_, uid, pseudo, _ := get_user_by_id(session.id_user)
		var keyuser = -1
		if get_key_allusers_by_id(session.id_user) == -1 {
			newid := get_new_id_user()
			us := &User{uid, session.id_user, newid, pseudo, false, UserSettings{true, true}, make(map[int]*IrcConnec), make(map[int]*Buffer), nil}
			all_users[newid] = us
			keyuser = newid
		} else {
			keyuser = get_key_allusers_by_id(session.id_user)
		}
		go all_users[keyuser].connect_server(session.server)
		log.Print("RESTORING : ", pseudo, " reconnecting on ", session.server)
	}

}

func restore_lost_channels(server string, server_id int, user_key int) {
	if serv_set.Restore_session == false {
		return
	}
	sessions := get_restore_sessions()
	log.Print("=== Crash ? Restoring channels....")

	for _, session := range sessions {
		if session.server == server {
			channels := strings.Split(session.channel, ",")
			for _, channel := range channels {
				if len(channel) > 1 {
					all_users[user_key].join_channel(server_id, channel)
				}
			}
			return
		}
	}

}
