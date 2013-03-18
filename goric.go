package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

func get_config_file() {
	content, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Panicln("conf.json error : ", err)
	}
	err = json.Unmarshal(content, &serv_set)
	if err != nil {
		log.Panicln("Error in conf.json : ", err)
	}
}

func main() {
	os.Chdir(os.Getenv("GOPATH") + "/src/goirc")
	rand.Seed(time.Now().UnixNano())
	get_config_file()
	if !serv_set.Show_log {
		log.SetOutput(ioutil.Discard)
	}
	connect_sql()
	get_preference()
	go restore_lost_server()
	start_http_server()
}
