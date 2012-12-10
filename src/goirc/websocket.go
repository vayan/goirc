package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/thoj/go-ircevent"
	"log"
	"strings"
)

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
		for pl, _ := range all_users {
			if all_users[pl].ws == ws {
				all_users[pl].ws = nil
				break
			}
		}
		fmt.Println(err)
	}
	log.Printf("recv :%s\n", buf)
	return buf, erri
}

func parsemsg(id_user int, msg string) {
	// if user connect 
	buff := strings.Split(msg, " ")
	switch buff[0] {
	case "/connect":
		go connect_server(buff[1], id_user)
		return
	case "/join":
		go join_channel(id_user, "irc.epitech.net:6667", buff[1])
		return
	case "/msg":
		go send_msg(id_user, "irc.epitech.net:6667", "#goirc", buff[1])
		return

	}

}

func WsHandle(ws *websocket.Conn) {
	log.Printf("Nouveau client %s\n", ws.Request().RemoteAddr)
	us := &User{"goricvayan", make(map[string]*irc.Connection), ws}
	all_users[1] = us
	for {

		if buff, err := ws_recv(ws); err != 1 {
			parsemsg(1, buff)
		} else {
			return
		}

	}
}
