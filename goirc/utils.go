package main

import (
	"bitbucket.org/vayan/gomin"
	"crypto/md5"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

//TODO : more flashy color and based on nick

func get_content_files(filename string, size int) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file_r := make([]byte, size)
	_, errs := file.Read(file_r)
	if errs != nil {
		log.Fatal(err)
	}
	return file_r
}

func write_files(filename string, content []byte) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write(content)
}

func optimize_static_files() {
	min_css := gomin.MinCSS(get_content_files("static/css/bootstrap.css", 50000))
	min_c := gomin.MinCSS(get_content_files("static/css/style.css", 50000))
	write_files("static/css/style-min.css", []byte(string(min_css)+string(min_c)))

	ircjs, err := gomin.MinJS(get_content_files("static/js/irc.js", 50000))
	if err != nil {
		log.Fatal(err)
	}
	write_files("static/js/irc-min.js", ircjs)
}

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
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Panicln("conf.json error : ", err)
	}
	err = json.Unmarshal(content, &serv_set)
	if err != nil {
		log.Panicln("Error in conf.json : ", err)
	}
}

func get_mail_hash(mail string) string {
	clean_mail := strings.Trim(strings.ToLower(mail), " ")
	h := md5.New()
	io.WriteString(h, clean_mail)
	return fmt.Sprintf("%x", h.Sum(nil))
}
