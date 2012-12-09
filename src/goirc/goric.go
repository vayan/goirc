package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	DB_NAME = "goirc"
	DB_USER = "goirc"
	DB_PASS = "rRfCKB6eMnDXNVZw"
	DB_SERV = "69.85.88.161"
)

type Page struct {
	Title string
}

func HandleErrorFatal(er error) bool {
	if er != nil {
		log.Fatal(er)
	}
	return false
}

func loadPage() *Page {
	title := "test"
	return &Page{Title: title}
}

func RenderHtml(w http.ResponseWriter, tmpl string, p *Page) {
	root_web := os.Getenv("GOIRC") + "/www/"
	t, _ := template.ParseFiles(root_web + tmpl + ".html")
	t.Execute(w, p)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage()
	RenderHtml(w, "index", p)
}

func main() {
	root_web := os.Getenv("GOIRC") + "/www/"
	if root_web == "" {
		log.Fatal("Set Root goric GOIRC plz")
	}

	fmt.Print("========Start connexion DB========\n\n")
	//username:password@hostspec/database
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASS+"@("+DB_SERV+":3306)/"+DB_NAME+"?charset=utf8")
	HandleErrorFatal(err)
	ar, err := db.Query("SELECT * FROM preference")
	HandleErrorFatal(err)
	fmt.Print(ar.Columns())

	fmt.Print("\n\n========Start goric web server========\n")
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(root_web+"css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(root_web+"img"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(root_web+"js"))))
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":1111", nil)
}
