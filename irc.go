package main

import (
	"github.com/thoj/go-ircevent"
	"log"
)

//ircobj := irc.IRC("goirctest"+strconv.Itoa(rand.Int()), "arheu")
//ircobj.VerboseCallbackHandler = true
// ircobj.Connect("irc.epitech.net:6667")
// ircobj.AddCallback("001", func(e *irc.Event) { ircobj.Join("#goirc") })
// ircobj.AddCallback("PRIVMSG", func(e *irc.Event) { ws_send(e.Nick+e.Message, sock_cli.ws) })

func send_msg(id_user int, server string, channel string, message string) {
	all_users[id_user].ircObj[server].irc.Privmsg(channel, message)
}

func join_channel(id_user int, server string, channel string) {
	all_users[id_user].ircObj[server].irc.Join(channel)
}

func connect_server(url string, id_user int) {
	con := irc.IRC(all_users[id_user].Nick, "arheu")
	all_users[id_user].ircObj[url] = &IrcConnec{con, ""}
	log.Print("creation obj irc")
	all_users[id_user].ircObj[url].irc.Connect(url)
	all_users[id_user].ircObj[url].irc.AddCallback("001", func(e *irc.Event) {
		ws_send("successserv]"+url, all_users[id_user].ws)
		all_users[id_user].ircObj[url].Nick = all_users[id_user].Nick
	})
	all_users[id_user].ircObj[url].irc.AddCallback("PRIVMSG", func(e *irc.Event) {
		ws_send(url+e.Arguments[0]+"]"+e.Message, all_users[id_user].ws)
	})
	all_users[id_user].ircObj[url].irc.AddCallback("366", func(e *irc.Event) {
		ws_send("successchan]"+url+"?"+e.Arguments[1], all_users[id_user].ws)
	})

}
