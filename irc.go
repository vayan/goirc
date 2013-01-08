package main

import (
	"strings"
)

func send_msg(id_user int, server int, message string) {
	user := all_users[id_user]

	user.ircObj[all_users[id_user].Buffers[server].id_serv].irc.Privmsg(all_users[id_user].Buffers[server].name, message)
}

func join_channel(id_user int, server int, channel string) {
	all_users[id_user].ircObj[server].irc.Join(channel)
}

func connect_server(url string, id_user int) {
	user := all_users[id_user]

	urlport := strings.Split(url, ":")
	if len(urlport) == 1 {
		url += ":6667"
	}
	id_buffer := user.get_new_id_buffer()
	user.add_buffer(urlport[0], url, id_buffer, id_buffer)
	user.add_connexion(user.Nick, "test", id_buffer)
	user.start_connexion(id_buffer, url)
	user.add_all_callback(id_buffer)
}
