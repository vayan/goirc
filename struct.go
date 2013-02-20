package main

import (
	"code.google.com/p/go.net/websocket"
	"container/list"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/thoj/go-ircevent"
)

/*
int of ircObj is the id of the buffer 
*/

type User struct {
	uid     string
	id      int
	Nick    string
	ircObj  map[int]*IrcConnec
	Buffers map[int]*Buffer
	ws      *websocket.Conn
}

type IrcConnec struct {
	irc  *irc.Connection
	Nick string
}

type BackLog struct {
	nick    string
	message string
	time    string
}

type Buffer struct {
	users   *list.List
	name    string
	addr    string
	id      int
	id_serv int //id du buffer network
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
	Data  map[string]interface{}
}

type RegisteringUser struct {
	InputMail      string
	InputPseudo    string
	InputPass      string
	InputPassVerif string
}
