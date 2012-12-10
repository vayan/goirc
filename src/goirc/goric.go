package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	connect_sql()
	get_preference()
	start_http_server()
}
