package main

import (
	"container/list"
	"log"
	"strconv"
	"strings"
)

//
// id_buffer (the key of the map) of the server buffer is the same as the ircObj key/id
//
//
// TODO : function send new buffer , for wraping all the ws_send
//

func (user *User) update_data_user() {
	_, id, pseudo, _ := get_user_by_uid(user.uid)

	user.id = id
	user.Nick = pseudo
}

// Retourne un ID pas utiliser pour buffer
func (user *User) get_new_id_buffer() int {
	if len(user.Buffers) == 0 {
		return 0
	}
	return len(user.Buffers) + 1
}

// Retourne ID buffer base sur son nom + id server
func (user *User) find_id_buffer(channel string, server int) int {
	channel = strings.ToLower(channel)
	for _, buff := range user.Buffers {
		if buff.name == channel && buff.id_serv == server {
			return buff.id
		}
	}
	//TODO : add gestion error
	return -1
}

// Retour l'id du serveur base sur l'id channel
func (user *User) find_server_by_channel(channel int) int {
	return user.Buffers[channel].id_serv
}

//TODO : verif duplicate buffer
func (user *User) add_buffer(name string, front_name string, addr string, id int, id_serv int) {
	// TODO : send new buffer to cl here, delete all other
	name = strings.ToLower(name)
	addr = strings.ToLower(addr)
	new_buffer := Buffer{list.New(), list.New(), name, front_name, addr, id, id_serv, false}
	user.Buffers[id] = &new_buffer

	//Restore friends
	server_buffer := user.Buffers[user.Buffers[id].id_serv]
	if session := user.raw_session[server_buffer.name]; session != nil {
		bff := strings.Split(session.friends, ",")
		for _, f := range bff {
			if len(f) > 0 {
				user.Buffers[id].friends.PushBack(f)
			}
		}
	}
}

func (user *User) add_con_loop(id_buffer int) {
	user.ircObj[id_buffer].irc.Loop()
}

func (user *User) start_connexion(id_buffer int, url string) {
	user.ircObj[id_buffer].irc.Connect(url)
}

func (user *User) send_all_buffer() {
	for _, buff := range user.Buffers {
		if buff.connected == true {
			ws_send("buffer]"+strconv.Itoa(buff.id)+"]"+buff.front_name+" "+user.ircObj[buff.id_serv].Nick, user.ws)
		}
	}
}

func (user *User) send_change_nick(id_buffer int, old_nick string, new_nick string) {
	for _, buff := range user.Buffers {
		if buff.id_serv == user.Buffers[id_buffer].id_serv {
			for e := user.Buffers[buff.id].users.Front(); e != nil; e = e.Next() {
				if e.Value.(ChannelUser).Nick == old_nick {
					val := user.Buffers[buff.id].users.Remove(e)
					chanuser := ChannelUser{new_nick, strings.Replace(new_nick, "@", "", 1), val.(ChannelUser).Color}
					user.Buffers[buff.id].users.PushBack(chanuser)
				}
			}
			ws_send("nick]"+strconv.Itoa(buff.id)+"]"+old_nick+" "+new_nick, user.ws)
		}
	}
}

func (user *User) close_buffer(id_buffer int) {
	if user.Buffers[id_buffer].id_serv == id_buffer {
		user.leave_network(id_buffer)
	} else {
		user.leave_channel(id_buffer)
	}
}

func (user *User) add_friend(id_buffer int, nick string) {
	//TODO : check nick is in buffer
	id_server := user.Buffers[id_buffer].id_serv
	go insert_new_friend_session(user.id, user.Buffers[id_server].name, nick)
	user.Buffers[id_buffer].friends.PushBack(nick)
}
