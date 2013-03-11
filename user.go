package main

import (
	"container/list"
	"github.com/thoj/go-ircevent"
	"html/template"
	"log"
	"strconv"
	"strings"
)

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
	new_buffer := Buffer{list.New(), name, front_name, addr, id, id_serv, false}
	user.Buffers[id] = &new_buffer
}

func (user *User) add_connexion(nick string, whois string, id_buffer int) {
	con := irc.IRC(nick, HOSTNAME_IRC)
	con.VerboseCallbackHandler = true //true for debug
	user.ircObj[id_buffer] = &IrcConnec{con, ""}
}

func (user *User) add_con_loop(id_buffer int) {
	user.ircObj[id_buffer].irc.Loop()
}

func (user *User) start_connexion(id_buffer int, url string) {
	user.ircObj[id_buffer].irc.Connect(url)
}

func (user *User) add_all_callback(id_buffer int) {
	user.on_connect(id_buffer)
	user.on_join(id_buffer)
	user.on_message(id_buffer)
	user.on_user_list(id_buffer)
	user.on_nick_used(id_buffer)
}

func (user *User) on_user_list(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("353", func(e *irc.Event) {
		id_buffer_chan := user.find_id_buffer(e.Arguments[2], id_buffer)
		if id_buffer_chan == -1 {
			id_buffer_chan = user.get_new_id_buffer()
			user.add_buffer(e.Arguments[2], e.Arguments[2], user.Buffers[id_buffer].addr+e.Arguments[2], id_buffer_chan, id_buffer)
		}
		arr := strings.Split(e.Message, " ")
		for _, val := range arr {
			//TODO : check other chara than @
			chanuser := ChannelUser{val, strings.Replace(val, "@", "", 1), GenerateColor()}
			user.Buffers[id_buffer_chan].users.PushBack(chanuser)
		}
	})
}

func (user *User) send_all_buffer() {
	for _, buff := range user.Buffers {
		if buff.connected == true {
			ws_send("buffer]"+strconv.Itoa(buff.id)+"]"+buff.front_name, user.ws)
		}
	}
}

func (user *User) on_connect(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("001", func(e *irc.Event) {
		user.Buffers[id_buffer].connected = true
		ws_send("buffer]"+strconv.Itoa(id_buffer)+"]"+user.Buffers[id_buffer].name, user.ws)
		log.Print("change nick")
		user.change_nick(id_buffer, user.Nick)
		log.Print("change nicked")
		ws_send("nick]"+strconv.Itoa(id_buffer)+"]"+user.ircObj[id_buffer].irc.GetNick(), user.ws)
		go insert_new_server_session(user.id, user.Buffers[id_buffer].name)
		user.ircObj[id_buffer].Nick = user.Nick
		go restore_lost_channels(user.Buffers[id_buffer].name, user.Buffers[id_buffer].id_serv, user.key)
	})
}

func (user *User) change_nick(id_buffer int, newnick string) {
	user.ircObj[user.Buffers[id_buffer].id_serv].irc.SendRaw("NICK " + newnick)
}

func (user *User) on_message(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("PRIVMSG", func(e *irc.Event) {
		buffer_name := e.Arguments[0]
		if buffer_name[0] != '#' {
			buffer_name = e.Nick
		}
		id_buffer_chan := user.find_id_buffer(buffer_name, id_buffer)
		if id_buffer_chan == -1 {
			id_buffer_chan = user.get_new_id_buffer()
			user.add_buffer(buffer_name, e.Nick, user.Buffers[id_buffer].addr, id_buffer_chan, id_buffer)
			ws_send("buffer]"+strconv.Itoa(id_buffer_chan)+"]"+e.Nick, user.ws)
		}
		log.Print(e.Arguments)
		go insert_new_message(user.id, user.Buffers[id_buffer].addr+e.Arguments[0], e.Nick, e.Message)
		go ws_send(strconv.Itoa(id_buffer_chan)+"]"+e.Nick+"]"+template.HTMLEscapeString(e.Message), user.ws)
	})
}

func (user *User) on_join(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("366", func(e *irc.Event) {
		id_buffer_chan := user.find_id_buffer(e.Arguments[1], id_buffer)
		if id_buffer_chan == -1 {
			id_buffer_chan := user.get_new_id_buffer()
			user.add_buffer(e.Arguments[1], e.Arguments[1], user.Buffers[id_buffer].addr+e.Arguments[1], id_buffer_chan, id_buffer)
		}
		user.Buffers[id_buffer_chan].connected = true
		ws_send("buffer]"+strconv.Itoa(id_buffer_chan)+"]"+e.Arguments[1], user.ws)
		insert_new_channel_session(user.id, user.Buffers[id_buffer].name, e.Arguments[1])
	})
}

func (user *User) on_nick_used(id_buffer int) {
	//TODO : randomize random pseudo
	user.ircObj[id_buffer].irc.AddCallback("433", func(e *irc.Event) {
		user.change_nick(id_buffer, "_"+user.Nick)
	})
}

func (user *User) close_buffer(id_buffer int) {
	if user.Buffers[id_buffer].id_serv == id_buffer {
		user.leave_network(id_buffer)
	} else {
		user.leave_channel(id_buffer)
	}
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
