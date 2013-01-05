package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"code.google.com/p/go.net/websocket"
	"github.com/thoj/go-ircevent"
)

type User struct {
	Nick   string
	ircObj map[string]*IrcConnec
	ws     *websocket.Conn
}

type IrcConnec struct {
	irc  *irc.Connection
	Nick string
}

type Preference struct {
	name              string
	descr             string
	short_descr       string
	long_descr        string
	base_url          string
	max_lenght_pseudo int
}

type Page struct {
	Title string
	Data  map[string]string
}

type RegisteringUser struct {
	InputMail      string
	InputPseudo    string
	InputPass      string
	InputPassVerif string
}
