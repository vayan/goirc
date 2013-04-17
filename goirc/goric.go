package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

//
// if go_path in conf.json is "env" the GOPATH env is used
//
//

func main() {
	rand.Seed(time.Now().UnixNano())
	get_config_file()
	if serv_set.Go_path == "env" {
		os.Chdir(os.Getenv("GOPATH") + "/goirc")
	} else if serv_set.Go_path == "" {
		log.Fatal("No path in config file")
	} else {
		os.Chdir(serv_set.Go_path)
	}
	if !serv_set.Show_log {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lshortfile)
	test_sql()
	get_preference()
	get_network()
	go restore_lost_server()
	optimize_static_files()
	start_http_server()
}
