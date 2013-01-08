package main

import (
	"github.com/thoj/go-ircevent"
	"log"
	"strconv"
)

func (user *User) get_new_id_buffer() int {
	if len(user.Buffers) == 0 {
		return 0
	}
	return len(user.Buffers) + 1
}

func (user *User) find_id_buffer(channel string, server int) int {
	for _, buff := range user.Buffers {
		if buff.name == channel && buff.id_serv == server {
			return buff.id
		}
	}
	return 0
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

func (user *User) on_connect(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("001", func(e *irc.Event) {
		ws_send("successserv]"+strconv.Itoa(id_buffer)+"]"+user.Buffers[id_buffer].name, user.ws)
		user.ircObj[id_buffer].Nick = user.Nick
	})
}

func (user *User) on_message(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("PRIVMSG", func(e *irc.Event) {
		id_buffer_chan := user.find_id_buffer(e.Arguments[0], id_buffer)
		go ws_send(strconv.Itoa(id_buffer_chan)+"]"+e.Message, user.ws)
	})
}

func (user *User) on_join(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("366", func(e *irc.Event) {
		id_buffer_chan := user.get_new_id_buffer()
		user.add_buffer(e.Arguments[1], user.Buffers[id_buffer].addr+e.Arguments[1], id_buffer_chan, id_buffer)
		ws_send("successchan]"+strconv.Itoa(id_buffer_chan)+"]"+e.Arguments[1], user.ws)
	})
}
