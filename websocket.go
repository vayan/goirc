package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
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

//parse les message client pour commande
func parsemsg(id_user int, msg string) {
	// if user is login 
	data := strings.SplitN((strings.SplitN(msg, "]", 2))[0], "?", 2)

	server := data[0]
	channel := data[1]
	buff := strings.Split((strings.SplitN(msg, "]", 2))[1], " ")

	switch buff[0] {
	case "/connect":
		log.Print("connexion server")
		go connect_server(buff[1], id_user)
		return
	case "/join":
		go join_channel(id_user, server, buff[1])
		return
	case "/msg":
		go send_msg(id_user, server, channel, buff[1])
		return

	}

}

func WsHandle(ws *websocket.Conn) {
	log.Printf("Nouveau client %s\n", ws.Request().RemoteAddr)
	us := &User{"goricvayan", make(map[string]*IrcConnec), ws}
	all_users[1] = us
	for {

		if buff, err := ws_recv(ws); err != 1 {
			parsemsg(1, buff)
		} else {
			return
		}

	}
}
