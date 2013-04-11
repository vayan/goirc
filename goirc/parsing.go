package main

import (
	"log"
	"strconv"
	"strings"
)

//TODO : epur str on all incoming message

func check_buffer_exist(id_buffer int, id_user int) bool {
	for id, _ := range all_users[id_user].Buffers {
		if id == id_buffer {
			return true
		}
	}
	return false
}

//parse les message client pour commande
func parsemsg(id_user int, msg string) {
	// TODO : Secure SECURE
	// TODO : TEST SECURE FFS
	user := all_users[id_user]
	data := strings.SplitN(msg, "]", 2)

	if data[0] == "co" {
		log.Print("User WS want to co")
		for pl, _ := range all_users {
			if all_users[pl].uid == data[1] {
				all_users[pl].ws = user.ws
				all_users[pl].online = true
				delete(all_users, id_user)
				log.Print("user find link etablish")
				all_users[pl].send_all_buffer()
				//TODO : generate new uid
				return
			}
		}
		log.Print("user not find create new instance")
		//TODO : check if uid exist in bdd
		user.online = true
		user.uid = data[1]
		user.update_data_user()
		user.ircObj = make(map[int]*IrcConnec)
		user.Buffers = make(map[int]*Buffer)
	} else if all_users[id_user].online == true {

		buffer_id, _ := strconv.Atoi(data[0])
		if len(data) <= 1 {
			return
		}
		buff_msg := data[1]
		buff := strings.Split(buff_msg, " ")

		switch buff[0] {
		case "/connect":
			if len(buff[1]) > 0 {
				go all_users[id_user].connect_server(buff[1])
			}
			return
		case "/join":
			if len(buff[1]) > 0 {
				go all_users[id_user].join_channel(user.find_server_by_channel(buffer_id), buff[1])
			}
			return
		case "/msg":
			if check_buffer_exist(buffer_id, id_user) && len(buff[1]) > 0 {
				go all_users[id_user].send_msg(buffer_id, buff[1])
			}
			return
		case "/close":
			if check_buffer_exist(buffer_id, id_user) {
				go all_users[id_user].close_buffer(buffer_id)
			}
			return
		case "/nick":
			if check_buffer_exist(buffer_id, id_user) && len(buff[1]) > 0 {
				go all_users[id_user].change_nick(buffer_id, buff[1])
			}
			return
		case "/me":
			if check_buffer_exist(buffer_id, id_user) && len(buff[1]) > 0 {
				go all_users[id_user].send_me(buffer_id, buff[1])
			}
			return
		case "/whois":
			if check_buffer_exist(buffer_id, id_user) && len(buff[1]) > 0 {
				go all_users[id_user].whois(buffer_id, buff[1])
			}
			return
		case "/addfriend":
			if check_buffer_exist(buffer_id, id_user) && len(buff[1]) > 0 {
				go all_users[id_user].add_friend(buffer_id, buff[1])
			}
			return
		default:
			if check_buffer_exist(buffer_id, id_user) {
				go all_users[id_user].send_msg(buffer_id, buff_msg)
			}
			return
		}
	}
}
