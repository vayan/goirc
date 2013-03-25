package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
)

//TODO : more flashy color
func GenerateColor() string {
	var r, g, b int

	r = (rand.Intn(256) + 60) / 2
	g = (rand.Intn(256) + 60) / 2
	b = (rand.Intn(256) + 60) / 2

	return ("rgb(" + strconv.Itoa(r) + ", " + strconv.Itoa(g) + ", " + strconv.Itoa(b) + ")")
}

func Atoi(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
	}
	return val
}

func generate_unique_uid(nick string) string {
	unique := string(rand.Int63()) + nick
	h := sha512.New()
	io.WriteString(h, unique)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ComparePassHash(pass string, hash string) bool {
	if EncryptPass(pass) == hash {
		return true
	}
	return false
}

func EncryptPass(pass string) string {
	h := sha512.New()

	io.WriteString(h, pass)

	return fmt.Sprintf("%x", h.Sum(nil))
}

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
