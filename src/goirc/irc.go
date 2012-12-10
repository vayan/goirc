package main

import (
	"github.com/thoj/go-ircevent"
)

//ircobj := irc.IRC("goirctest"+strconv.Itoa(rand.Int()), "arheu")
//ircobj.VerboseCallbackHandler = true
// ircobj.Connect("irc.epitech.net:6667")
// ircobj.AddCallback("001", func(e *irc.Event) { ircobj.Join("#goirc") })
// ircobj.AddCallback("PRIVMSG", func(e *irc.Event) { ws_send(e.Nick+e.Message, sock_cli.ws) })

func send_msg(id_user int, server string, channel string, message string) {
	all_users[id_user].ircObj[server].Privmsg(channel, message)
}

func join_channel(id_user int, server string, channel string) {
	all_users[id_user].ircObj[server].Join(channel)
	all_users[id_user].ircObj[server].AddCallback("366", func(e *irc.Event) { ws_send("connect to channel !!", all_users[id_user].ws) })
}

func connect_server(url string, id_user int) {
	all_users[id_user].ircObj[url] = irc.IRC(all_users[id_user].Nick, "arheu")
	all_users[id_user].ircObj[url].Connect(url)
	all_users[id_user].ircObj[url].AddCallback("001", func(e *irc.Event) { ws_send("connect !!", all_users[id_user].ws) })

}
