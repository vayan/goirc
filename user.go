package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/thoj/go-ircevent"
	"log"
	"strconv"
)

// Retourne un ID pas utiliser pour la list d'usr 

func get_new_id_user() int {
	if len(all_users) == 0 {
		return 0
	}
	return len(all_users) + 1
}

//get client by id 

func get_user_id(id int) *User {
	for _, us := range all_users {
		if us.id == id {
			return us
		}
	}
	return nil
}

// get id client by ws 
func get_user_ws(ws *websocket.Conn) int {
	for pl, _ := range all_users {
		if all_users[pl].ws == ws {
			return pl
		}
	}
	return -1
}

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

// Retourne ID buffer base sur son nom + server
func (user *User) find_id_buffer(channel string, server int) int {
	for _, buff := range user.Buffers {
		if buff.name == channel && buff.id_serv == server {
			return buff.id
		}
	}
	return 0
}

// Retour l'id du serveur base sur l'id channel
func (user *User) find_server_by_channel(channel int) int {
	return user.Buffers[channel].id_serv
}

func (user *User) add_buffer(name string, addr string, id int, id_serv int) {
	new_buffer := Buffer{name, addr, id, id_serv}
	user.Buffers[id] = &new_buffer
}

func (user *User) add_connexion(nick string, whois string, id_buffer int) {
	con := irc.IRC(user.Nick, "arheu")
	user.ircObj[id_buffer] = &IrcConnec{con, ""}
}

func (user *User) start_connexion(id_buffer int, url string) {
	user.ircObj[id_buffer].irc.Connect(url)
}

func (user *User) add_all_callback(id_buffer int) {
	user.on_connect(id_buffer)
	user.on_join(id_buffer)
	user.on_message(id_buffer)
}

func (user *User) send_all_buffer() {
	for _, buff := range user.Buffers {
		ws_send("buffer]"+strconv.Itoa(buff.id)+"]"+buff.name, user.ws)
	}
}

func (user *User) on_connect(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("001", func(e *irc.Event) {
		ws_send("buffer]"+strconv.Itoa(id_buffer)+"]"+user.Buffers[id_buffer].name, user.ws)
		user.ircObj[id_buffer].Nick = user.Nick
	})
}

func (user *User) on_message(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("PRIVMSG", func(e *irc.Event) {
		id_buffer_chan := user.find_id_buffer(e.Arguments[0], id_buffer)
		log.Print(e.Arguments)
		go insert_new_message(user.id, user.Buffers[id_buffer].addr+e.Arguments[0], e.Nick, e.Message)
		go ws_send(strconv.Itoa(id_buffer_chan)+"]"+e.Nick+"]"+e.Message, user.ws)
	})
}

func (user *User) on_join(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("366", func(e *irc.Event) {
		id_buffer_chan := user.get_new_id_buffer()
		user.add_buffer(e.Arguments[1], user.Buffers[id_buffer].addr+e.Arguments[1], id_buffer_chan, id_buffer)
		ws_send("buffer]"+strconv.Itoa(id_buffer_chan)+"]"+e.Arguments[1], user.ws)
	})
}
