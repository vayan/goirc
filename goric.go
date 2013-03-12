package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

//TODO : gerer reconnexion server/channel si crash

func get_config_file() {
	log.Print("=== Get config from conf.json")
	content, err := ioutil.ReadFile("conf.json")
	if err != nil {
		//TODO : error
	}
	err = json.Unmarshal(content, &set)
	if err != nil {
		log.Print("Error in conf.json : ", err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	get_config_file()
	connect_sql()
	get_preference()
	go restore_lost_server()
	start_http_server()
}
