package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
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

func WsHandle(ws *websocket.Conn) {
	log.Printf("Nouveau client %s\n", ws.Request().RemoteAddr)
	us := &User{"goricvayan", make(map[int]*IrcConnec), make(map[int]*Buffer), ws}
	all_users[1] = us
	for {

		if buff, err := ws_recv(ws); err != 1 {
			go parsemsg(1, buff)
		} else {
			return
		}

	}
}
