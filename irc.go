package main

import (
	"github.com/thoj/go-ircevent"
	"html/template"
	"strings"
)

func (user *User) add_connexion(nick string, whois string, id_buffer int) {
	con := irc.IRC(nick, serv_set.Hostname_irc)
	con.VerboseCallbackHandler = false //true for debug
	user.ircObj[id_buffer] = &IrcConnec{con, ""}
}

func (user *User) send_msg(server int, message string) {
	msg := template.HTMLEscapeString(message)

	user.ircObj[user.Buffers[server].id_serv].irc.Privmsg(user.Buffers[server].name, msg)
	go insert_new_message(user.id, user.Buffers[server].addr, user.ircObj[user.Buffers[server].id_serv].irc.GetNick(), msg)
}

func (user *User) join_channel(server int, channel string) {
	user.ircObj[server].irc.Join(channel)
}

func (user *User) connect_server(url string) {
	urlport := strings.Split(url, ":")
	if len(urlport) == 1 {
		url += ":6667"
	}
	id_buffer := user.get_new_id_buffer()
	user.add_buffer(urlport[0], urlport[0], url, id_buffer, id_buffer)
	user.add_connexion(user.Nick, "test", id_buffer)
	user.start_connexion(id_buffer, url)
	user.add_all_callback(id_buffer)
	user.add_con_loop(id_buffer)
}

func (user *User) leave_channel(id_buffer_chan int) {
	id_ircobj := user.Buffers[id_buffer_chan].id_serv
	user.ircObj[id_ircobj].irc.Part(user.Buffers[id_buffer_chan].name)
	delete(user.Buffers, id_buffer_chan)
	//TODO : remove in DB
}

func (user *User) leave_network(id_buffer_chan int) {
	id_ircobj := user.Buffers[id_buffer_chan].id_serv
	user.ircObj[id_ircobj].irc.Quit()

	for key, buff := range user.Buffers {
		if buff.id_serv == id_ircobj {
			delete(user.Buffers, key)
		}
	}
	delete(user.ircObj, id_ircobj)
	//TODO : Remove in db
}

func (user *User) change_nick(id_buffer int, newnick string) {
	user.ircObj[user.Buffers[id_buffer].id_serv].irc.Nick(newnick)
}
