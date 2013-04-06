package main

import (
	"github.com/thoj/go-ircevent"
	"html/template"
	"log"
	"strconv"
	"strings"
)

func (user *User) add_all_callback(id_buffer int) {
	user.on_connect(id_buffer)
	user.on_me_join(id_buffer)
	user.on_message(id_buffer)
	user.on_user_list(id_buffer)
	user.on_nick_used(id_buffer)
	user.on_nick_change(id_buffer)
	user.on_part(id_buffer)
	user.on_join(id_buffer)
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

func (user *User) on_connect(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("001", func(e *irc.Event) {
		user.Buffers[id_buffer].connected = true
		ws_send("buffer]"+strconv.Itoa(id_buffer)+"]"+user.Buffers[id_buffer].name+" "+user.ircObj[id_buffer].Nick, user.ws)
		user.ircObj[id_buffer].Nick = user.ircObj[id_buffer].irc.GetNick()
		ws_send("nick]"+strconv.Itoa(id_buffer)+"]"+user.ircObj[id_buffer].irc.GetNick(), user.ws)
		user.ircObj[id_buffer].Nick = user.Nick
		go insert_new_server_session(user.id, user.Buffers[id_buffer].name)
		go restore_lost_channels(user.Buffers[id_buffer].name, user.Buffers[id_buffer].id_serv, user.key)
	})
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
			go ws_send("buffer]"+strconv.Itoa(id_buffer_chan)+"]"+e.Nick+" "+user.ircObj[id_buffer].Nick, user.ws)
		}
		log.Print(e.Arguments)
		go insert_new_message(user.id, user.Buffers[id_buffer].addr+e.Arguments[0], e.Nick, e.Message)
		go ws_send(strconv.Itoa(id_buffer_chan)+"]"+e.Nick+"]"+template.HTMLEscapeString(e.Message), user.ws)
	})
}

func (user *User) on_me_join(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("366", func(e *irc.Event) {
		id_buffer_chan := user.find_id_buffer(e.Arguments[1], id_buffer)
		if id_buffer_chan == -1 {
			id_buffer_chan := user.get_new_id_buffer()
			user.add_buffer(e.Arguments[1], e.Arguments[1], user.Buffers[id_buffer].addr+e.Arguments[1], id_buffer_chan, id_buffer)
		}
		user.Buffers[id_buffer_chan].connected = true
		go ws_send("buffer]"+strconv.Itoa(id_buffer_chan)+"]"+e.Arguments[1]+" "+user.ircObj[id_buffer].Nick, user.ws)
		go insert_new_channel_session(user.id, user.Buffers[id_buffer].name, e.Arguments[1])
	})
}

func (user *User) on_part(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("PART", func(e *irc.Event) {
		id_buffer_chan := user.find_id_buffer(e.Arguments[0], id_buffer)
		if id_buffer_chan == -1 {
			return
		}
		go ws_send("part]"+strconv.Itoa(id_buffer_chan)+"]"+e.Nick, user.ws)
		for j := user.Buffers[id_buffer_chan].users.Front(); j != nil; j = j.Next() {
			if j.Value.(ChannelUser).Nick == e.Nick {
				user.Buffers[id_buffer_chan].users.Remove(j)
				return
			}
		}

	})
}

func (user *User) on_join(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("JOIN", func(e *irc.Event) {
		var id_buffer_chan int
		if id_buffer_chan = user.find_id_buffer(e.Message, id_buffer); id_buffer_chan == -1 {
			return
		}
		log.Print("join user add to list")
		chanuser := ChannelUser{e.Nick, strings.Replace(e.Nick, "@", "", 1), GenerateColor()}
		user.Buffers[id_buffer_chan].users.PushBack(chanuser)
		go ws_send("join]"+strconv.Itoa(id_buffer_chan)+"]"+e.Nick, user.ws)
	})
}

func (user *User) on_nick_change(id_buffer int) {
	user.ircObj[id_buffer].irc.AddCallback("NICK", func(e *irc.Event) {
		if user.ircObj[user.Buffers[id_buffer].id_serv].Nick == e.Nick {
			user.ircObj[user.Buffers[id_buffer].id_serv].Nick = e.Message
		}
		go user.send_change_nick(id_buffer, e.Nick, e.Message)
	})
}

func (user *User) on_nick_used(id_buffer int) {
	//TODO : randomize random pseudo
	user.ircObj[id_buffer].irc.AddCallback("433", func(e *irc.Event) {
		user.change_nick(id_buffer, "_"+user.Nick)
	})
}
