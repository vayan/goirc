package main

import (
	//"code.google.com/p/go-mysql-driver/mysql"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

const (
	DB_NAME = "goirc"
	DB_USER = "goirc"
	DB_PASS = "rRfCKB6eMnDXNVZw"
)

type Page struct {
	Title string
}

func loadPage() *Page {
	title := "test"
	return &Page{Title: title}
}

func RenderHtml(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPage()
	RenderHtml(w, "index", p)
}

func main() {
	db, err := sql.Open("mysql", "goirc:rRfCKB6eMnDXNVZw@tcp(vayan.fr:5555)/goirc?charset=utf8")
	if err {
		fmt.Print("error")
	}
	tx, _ := db.Begin()

	row, _ := tx.Query("SELECT * FROM preference")

	fmt.Print(row)

	fmt.Print("Start goric web server\n")
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":1111", nil)
}
