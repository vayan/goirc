package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"net/http"
	"strconv"
)

func GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	if need_perm(REGIST, r) {
		// TODO : stats
	}
}

func GetFriendsHandler(w http.ResponseWriter, r *http.Request) {
	if need_perm(REGIST, r) {
		session, _ := store.Get(r, serv_set.Cookie_session)
		if us := get_user_id(session.Values["id"].(int)); us != nil {
			idbuffer := Atoi(r.FormValue("idbuffer"))
			if val, ok := us.Buffers[idbuffer]; ok && us.Buffers[idbuffer].friends.Len() > 0 {
				jsonres := "{ \"FriendList\":["
				for e := val.friends.Front(); e != nil; e = e.Next() {
					b, _ := json.Marshal(e.Value.(string))
					jsonres += string(b) + ","
				}
				jsonres = jsonres[:len(jsonres)-1]
				jsonres += "]}"
				fmt.Fprint(w, jsonres)
			}
		}
	}
}

func SetChanHandler(w http.ResponseWriter, r *http.Request) {
	if need_perm(REGIST, r) {
		//TODO : test si ws active
		var allserv = make(map[string]string)
		session, _ := store.Get(r, serv_set.Cookie_session)

		if us := get_user_id(session.Values["id"].(int)); us != nil {
			for _, irco := range us.Buffers {
				if irco.id == irco.id_serv && irco.connected == true {
					allserv[irco.name] = strconv.Itoa(irco.id)
				}
			}
			b, _ := json.Marshal(allserv)
			fmt.Fprint(w, string(b))
		}
	}
}

func ActionBacklogHandler(w http.ResponseWriter, r *http.Request) {
	//TODO : check id correct
	idbuffer := Atoi(r.FormValue("idbuffer"))
	session, _ := store.Get(r, serv_set.Cookie_session)

	if need_perm(REGIST, r) {
		if user := get_user_id(session.Values["id"].(int)); user != nil {
			buffers := user.Buffers
			for _, buff := range buffers {
				if buff.id == idbuffer {
					backlog := get_backlog(user.id, user.Buffers[idbuffer].addr)
					b, _ := json.Marshal(backlog)
					fmt.Fprint(w, string(b))
					return
				}
			}
		}
	}
}

func UsersListHandler(w http.ResponseWriter, r *http.Request) {
	if need_perm(REGIST, r) {
		jsonres := "{ \"UserList\":["
		id := Atoi(r.FormValue("channel"))
		session, _ := store.Get(r, serv_set.Cookie_session)
		if us := get_user_id(session.Values["id"].(int)); us != nil {

			if _, ok := us.Buffers[id]; ok && us.Buffers[id].users.Len() > 0 {
				for e := us.Buffers[id].users.Front(); e != nil; e = e.Next() {
					b, _ := json.Marshal(e.Value.(ChannelUser))
					jsonres += string(b) + ","
				}
				jsonres = jsonres[:len(jsonres)-1]
				jsonres += "]}"
				fmt.Fprint(w, jsonres)
			}
		}
	}
}

func GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, serv_set.Cookie_session)
	if user := get_user_id(session.Values["id"].(int)); need_perm(REGIST, r) && user != nil {
		jsonres, _ := json.Marshal(user.Settings)
		fmt.Fprint(w, string(jsonres))
	}
}
