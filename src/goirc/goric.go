package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"code.google.com/p/go.net/websocket"
	"github.com/thoj/go-ircevent"
	"math/rand"
	"time"
)

var all_users = make(map[int]*User)

type User struct {
	Nick   string
	ircObj map[string]*irc.Connection
	ws     *websocket.Conn
}

func main() {
	rand.Seed(time.Now().UnixNano())
	connect_sql()
	get_preference()
	start_http_server()
}
