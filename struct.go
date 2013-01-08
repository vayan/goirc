package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"code.google.com/p/go.net/websocket"
	"github.com/thoj/go-ircevent"
)

/*
int of ircObj is the id of the buffer 
*/

type User struct {
	Nick    string
	ircObj  map[int]*IrcConnec
	Buffers map[int]*Buffer
	ws      *websocket.Conn
}

type IrcConnec struct {
	irc  *irc.Connection
	Nick string
}

type Buffer struct {
	name    string
	addr    string
	id      int
	id_serv int
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
