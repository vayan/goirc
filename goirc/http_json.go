package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"html/template"
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
		us := get_user_id(session.Values["id"].(int))
		idbuffer := Atoi(r.FormValue("idbuffer"))
		if val, ok := us.Buffers[idbuffer]; ok {
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

func SetChanHandler(w http.ResponseWriter, r *http.Request) {
	if need_perm(REGIST, r) {
		//TODO : test si ws active
		//TODO : Json

		var allserv string
		session, _ := store.Get(r, serv_set.Cookie_session)
		us := get_user_id(session.Values["id"].(int))

		for _, irco := range us.Buffers {
			if irco.name[0] != '#' {
				allserv += "<option value='" + strconv.Itoa(irco.id) + "'>" + irco.name + "</option>"
			}
		}
		p := &Page{
			Title: "IRC in your browser",
			Data:  map[string]interface{}{"servers": template.HTML(allserv)}}
		RenderHtml(w, "ajx/set-channels", p)
	}
}

func ActionBacklogHandler(w http.ResponseWriter, r *http.Request) {
	//TODO : JSON all that
	//TODO : check id correct
	idbuffer := Atoi(r.FormValue("idbuffer"))
	session, _ := store.Get(r, serv_set.Cookie_session)

	if need_perm(REGIST, r) {
		user := get_user_id(session.Values["id"].(int))
		buffers := user.Buffers
		for _, buff := range buffers {
			if buff.id == idbuffer {
				backlog := get_backlog(user.id, user.Buffers[idbuffer].addr)
				for _, log := range backlog {
					fmt.Fprint(w, "<tr class=\"msg\"><td class=\"pseudo "+log.nick+"\">"+log.nick+"</td><td class=\"message\"><div class='messagediv'>"+log.message+"</div></td><td class=\"time\">"+log.time+"</td></tr>")
				}
				return
			}
		}

	}
}

func UsersListHandler(w http.ResponseWriter, r *http.Request) {
	if need_perm(REGIST, r) {
		jsonres := "{ \"UserList\":["
		id := Atoi(r.FormValue("channel"))
		session, _ := store.Get(r, serv_set.Cookie_session)
		us := get_user_id(session.Values["id"].(int))

		if _, ok := us.Buffers[id]; ok {
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

func GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, serv_set.Cookie_session)
	if need_perm(REGIST, r) {
		user := get_user_id(session.Values["id"].(int))
		jsonres, _ := json.Marshal(user.Settings)
		fmt.Fprint(w, string(jsonres))
	}
}
