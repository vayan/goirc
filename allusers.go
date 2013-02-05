package main

import (
	"code.google.com/p/go.net/websocket"
)

// Retourne un ID pas utiliser pour la list d'usr 
func get_new_id_user() int {
	if len(all_users) == 0 {
		return 0
	}
	return len(all_users) + 1
}

//get client by id 
func get_user_id(id int) *User {
	for _, us := range all_users {
		if us.id == id {
			return us
		}
	}
	return nil
}

// get id client by ws 
func get_user_ws(ws *websocket.Conn) int {
	for pl, _ := range all_users {
		if all_users[pl].ws == ws {
			return pl
		}
	}
	return -1
}
