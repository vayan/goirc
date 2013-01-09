package main

import (
	"strconv"
	"strings"
)

//parse les message client pour commande
func parsemsg(id_user int, msg string) {
	// if user is login 
	user := all_users[id_user]
	data := strings.Split(msg, "]")

	buffer_id, _ := strconv.Atoi(data[0])
	buff_msg := data[1]

	buff := strings.Split(buff_msg, " ")

	switch buff[0] {
	case "/connect":
		go connect_server(buff[1], id_user)
		return
	case "/join":
		go join_channel(id_user, user.find_server_by_channel(buffer_id), buff[1])
		return
	case "/msg":
		go send_msg(id_user, buffer_id, buff[1])
		return
	default:
		go send_msg(id_user, buffer_id, buff_msg)
		return
	}
}
