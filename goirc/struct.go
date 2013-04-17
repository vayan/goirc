package main

import (
	"code.google.com/p/go.net/websocket"
	"container/list"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thoj/go-ircevent"
)

/*
int of ircObj is the id of the buffer
*/

type User struct {
	uid         string
	id          int
	key         int
	Nick        string
	online      bool
	connecting  bool
	Settings    UserSettings
	ircObj      map[int]*IrcConnec
	Buffers     map[int]*Buffer
	ws          *websocket.Conn
	raw_session map[string]*RestoreSession
}

type UserSettings struct {
	Notify       bool
	Save_session bool
}

type ChannelUser struct {
	Nick      string
	NickClean string
	Color     string
}

type IrcConnec struct {
	irc  *irc.Connection
	Nick string
}

type BackLog struct {
	Nick    string
	Message string
	Time    string
}

type Buffer struct {
	users      *list.List
	friends    *list.List
	name       string
	front_name string
	addr       string
	id         int
	id_serv    int //id du buffer network or ircobj
	connected  bool
}

type Server_Settings struct {
	Root_web        string
	Go_path         string
	Port_http       string
	Cookie_session  string
	DB_server       string
	DB_name         string
	DB_user         string
	DB_pass         string
	Restore_session bool
	Show_log        bool
	Show_log_irc    bool
	Log_in_file     bool
	Hostname_irc    string
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

type RestoreSession struct {
	id_user int
	server  string
	channel string
	friends string
}

type UserStats struct {
	ConnectedServer int
	ConnctedChan    int
	CloudWords      string
	MsgSend         int
}

type Network struct {
	id                int
	name              string
	adress            string
	port              int
	limit             int
	current_connected int
}
