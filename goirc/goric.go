package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

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
	optimize_static_files()
	start_http_server()
}
