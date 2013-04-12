package main

import (
	"github.com/thoj/go-ircevent"
	"strconv"
	"strings"
)

func (user *User) add_connexion(nick string, whois string, id_buffer int) {
	con := irc.IRC(nick, serv_set.Hostname_irc)
	con.VerboseCallbackHandler = true //true for debug
	if user.ircObj != nil {
		user.ircObj[id_buffer] = &IrcConnec{con, ""}
	}

}

func (user *User) send_msg(server int, message string) {
	if buf, ok := user.Buffers[server]; ok && buf.connected == true {
		if co, ok := user.ircObj[buf.id_serv]; ok && co.irc != nil {
			co.irc.Privmsg(buf.name, message)
			go insert_new_message(user.id, buf.addr, co.irc.GetNick(), message)
		}
	}
}

func (user *User) join_channel(server int, channel string) {
	if buf, ok := user.Buffers[server]; ok && buf.connected == true {
		user.ircObj[server].irc.Join(channel)
	}
}

func (user *User) connect_server(url string) {
	urlport := strings.Split(url, ":")
	if len(urlport) == 1 {
		url += ":6667"
	}
	if url != ":6667" {
		for _, buf := range user.Buffers {
			if buf.name == urlport[0] {
				return
			}
		}
		id_buffer := user.get_new_id_buffer()
		if id_buffer != -1 {
			user.connecting = true
			user.add_buffer(urlport[0], urlport[0], url, id_buffer, id_buffer)
			user.add_connexion(user.Nick, "test", id_buffer)
			user.start_connexion(id_buffer, url)
			user.add_all_callback(id_buffer)
			user.connecting = false
			user.add_con_loop(id_buffer)
		}
	}
}

func (user *User) leave_channel(id_buffer_chan int, remove_session bool) {
	buff, ok := user.Buffers[id_buffer_chan]
	if !(ok) || (buff.id == buff.id_serv) {
		return
	}
	id_ircobj := buff.id_serv
	if buff.connected == true {
		if remove_session {
			go remove_channel_session(user.id, user.Buffers[id_ircobj].name, user.Buffers[id_buffer_chan].name)
		} else {
			go ws_send("leave]"+strconv.Itoa(id_buffer_chan), user.ws)
		}
		user.ircObj[id_ircobj].irc.Part(user.Buffers[id_buffer_chan].name)
		delete(user.Buffers, id_buffer_chan)
	}
}

func (user *User) leave_network(id_buffer_chan int) {
	//TODO : remove all buffer channel
	id_ircobj := user.Buffers[id_buffer_chan].id_serv
	if co, ok := user.Buffers[id_ircobj]; ok {
		co.connected = false
		remove_server_session(user.id, co.name)
		user.ircObj[id_ircobj].irc.Quit()

		for key, buff := range user.Buffers {
			if buff.id_serv == id_ircobj {
				delete(user.Buffers, key)
			}
		}
		delete(user.ircObj, id_ircobj)
	}
}

func (user *User) whois(id_buffer int, nick string) {
	if buf, ok := user.Buffers[id_buffer]; ok && buf.connected == true {
		id_server := user.Buffers[id_buffer].id_serv
		user.ircObj[id_server].irc.SendRawf("WHOIS %s %s", user.Buffers[id_server].name, nick)
	}
}

func (user *User) send_me(id_buffer int, msg string) {
	if buf, ok := user.Buffers[id_buffer]; ok && buf.connected == true {
		user.ircObj[user.Buffers[id_buffer].id_serv].irc.SendRawf("PRIVMSG %s :\x01ACTION %s\x01", user.Buffers[id_buffer].name, msg)
	}
}

func (user *User) change_nick(id_buffer int, newnick string) {
	if buf, ok := user.Buffers[id_buffer]; ok && buf.connected == true {
		user.ircObj[user.Buffers[id_buffer].id_serv].irc.Nick(newnick)
	}
}
