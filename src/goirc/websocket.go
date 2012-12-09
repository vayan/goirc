package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/thoj/go-ircevent"
	"log"
	"math/rand"
	"strconv"
)

var client = make(map[Connection]int)

type Connection struct {
	ws       *websocket.Conn
	clientip string
	irc      *irc.Connection
}

func ws_send(buf string, ws *websocket.Conn) {
	if err := websocket.Message.Send(ws, buf); err != nil {
		log.Println(err)
	}
	log.Printf("send:%s\n", buf)
}

func ws_recv(ws *websocket.Conn) (string, int) {
	var buf string
	erri := 0

	if err := websocket.Message.Receive(ws, &buf); err != nil {
		erri = 1
		for pl, _ := range client {
			if pl.ws == ws {
				delete(client, pl)
				break
			}
		}
		fmt.Println(err)
	}
	log.Printf("recv :%s\n", buf)
	return buf, erri
}

func WsHandle(ws *websocket.Conn) {
	ircobj := irc.IRC("goirctest"+strconv.Itoa(rand.Int()), "arheu")
	//ircobj.VerboseCallbackHandler = true
	sock_cli := Connection{ws, ws.Request().RemoteAddr, ircobj}
	ircobj.Connect("irc.epitech.net:6667")
	ircobj.AddCallback("001", func(e *irc.Event) { ircobj.Join("#goirc") })
	ircobj.AddCallback("PRIVMSG", func(e *irc.Event) { ws_send(e.Nick+e.Message, sock_cli.ws) })
	fmt.Printf("\nNouveau client %s\n", sock_cli.clientip)
	client[sock_cli] = 0
	for {

		if _, err := ws_recv(ws); err == 1 {
			return
		}

	}
}
