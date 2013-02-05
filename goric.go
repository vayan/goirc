package main

import (
	"math/rand"
	"time"
)

//TODO : gerer reconnexion server / channel si crash

func main() {
	rand.Seed(time.Now().UnixNano())
	connect_sql()
	get_preference()
	start_http_server()
}
