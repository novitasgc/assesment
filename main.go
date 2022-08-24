package main

import (
	"database/sql"
	"net/http"

	"fmt"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var assesment = template.Must(template.ParseGlob("C:/Users/MSI-PC/Desktop/nop/assesment/*.html"))

func main() {
	http.HandleFunc("/", Record)
	http.HandleFunc("/add.html", Add)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/edited", Edited)
	http.HandleFunc("/delete", Delete)
	fmt.Println("Server Continue")
	http.ListenAndServe(":8080", nil)

}

func sqlconnection() (conection *sql.DB) {

	conection, err := sql.Open("mysql", "root:@tcp(127.0.0.1)/assesment")
	if err != nil {
		panic(err.Error())
	}
	return conection
}

type tabledata struct {
	Id     int
	Title  string
	Genre  string
	Rating int
}

func Record(w http.ResponseWriter, r *http.Request) {
	connectionEstablish := sqlconnection()
	tables, err := connectionEstablish.Query("SELECT * FROM records")

	if err != nil {
		panic(err.Error())
	}

	data := tabledata{}
	arrange := []tabledata{}

	for tables.Next() {
		var id int
		var title, genre string
		var rating int
		err = tables.Scan(&id, &title, &genre, &rating)

		if err != nil {
			panic(err.Error())
		}
		data.Id = id
		data.Title = title
		data.Genre = genre
		data.Rating = rating

		arrange = append(arrange, data)
	}
	assesment.ExecuteTemplate(w, "record", arrange)
}

func Add(w http.ResponseWriter, r *http.Request) {
	assesment.ExecuteTemplate(w, "add", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		genre := r.FormValue("genre")
		rating := r.FormValue("rating")

		connectionEstablish := sqlconnection()
		insert, err := connectionEstablish.Prepare("INSERT INTO records(title, genre, rating) VALUES (?,?,?)")

		if err != nil {
			panic(err.Error())
		}
		insert.Exec(title, genre, rating)
		http.Redirect(w, r, "/", 301)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idtable := r.URL.Query().Get("id")
	connectionEstablish := sqlconnection()
	delete, err := connectionEstablish.Prepare("DELETE FROM records WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	delete.Exec(idtable)
	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	idtable := r.URL.Query().Get("id")
	connectionEstablish := sqlconnection()
	tables, err := connectionEstablish.Query("SELECT * FROM records WHERE id=?", idtable)

	if err != nil {
		panic(err.Error())
	}

	data := tabledata{}

	for tables.Next() {
		var id int
		var title string
		var genre string
		var rating int
		err = tables.Scan(&id, &title, &genre, &rating)

		if err != nil {
			panic(err.Error())
		}
		data.Id = id
		data.Title = title
		data.Genre = genre
		data.Rating = rating
		fmt.Println(data)
	}
	assesment.ExecuteTemplate(w, "edit", data)
}

func Edited(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		Id := r.FormValue("id")
		Title := r.FormValue("title")
		Genre := r.FormValue("genre")
		Rating := r.FormValue("rating")

		connectionEstablish := sqlconnection()
		edit, err := connectionEstablish.Prepare("UPDATE records SET title=?, genre=?, rating=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		edit.Exec(Title, Genre, Rating, Id)
		http.Redirect(w, r, "/", 301)
	}
}
