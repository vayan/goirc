package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
)

func ws_send(buf string, ws *websocket.Conn) {
	if ws == nil || len(buf) == 0 {
		return
	}
	if err := websocket.Message.Send(ws, buf); err != nil {
		log.Println(err)
	}
	log.Printf("send : '%s'", buf)
}

func ws_recv(ws *websocket.Conn) (string, int) {
	var buf string
	erri := 0

	if err := websocket.Message.Receive(ws, &buf); err != nil {
		erri = 1
		for pl, _ := range all_users {
			if all_users[pl].ws == ws {
				all_users[pl].ws = nil
				all_users[pl].online = false
				break
			}
		}
		fmt.Println(err)
	}
	if len(buf) < 1 {
		return buf, 1
	}
	if len(buf) > 200 {
		return buf[0:200], erri
	}
	log.Printf("recv : '%s'", buf)
	return buf, erri
}

func WsHandle(ws *websocket.Conn) {
	log.Printf("Nouveau client %s", ws.Request().RemoteAddr)
	newid := get_new_id_user()
	us := &User{"nil", 0, newid, "Anon3123123123", false, UserSettings{true, true}, nil, nil, ws}
	all_users[newid] = us
	for {

		if buff, err := ws_recv(ws); err != 1 {
			go parsemsg(get_user_ws(ws), buff)
		} else {
			return
		}

	}
}
