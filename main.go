package main

import (
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"log"
	"strconv"
)

type Post struct {
	Id int
	Title string
	Body string
}

var Db *sql.DB
var templates = template.Must(template.New("t").ParseGlob("views/**/*.html"))

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=postgres dbname=gwp password=gwp sslmode=disable")
	
	if err != nil {
		panic(err)
	}
}

func Posts() (posts []Post, err error) {
	rows, err := Db.Query("select id, title, body from posts")

	if err != nil {
		return
	}

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Body)

		if err != nil {
			return
		}

		posts = append(posts, post)
	}

	rows.Close()
	return
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "posts/new", "") 
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	post := Post{}
	title := r.FormValue("title")
	body := r.FormValue("body")
	statement := "insert into posts (title, body) values ($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(title, body).Scan(&post.Id)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusMovedPermanently)
}

func PostIndex(w http.ResponseWriter, r *http.Request) {
	posts, _ := Posts()
	templates.ExecuteTemplate(w, "posts/index", posts) 
}

func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("select id, title, body from posts where id = $1", id).Scan(&post.Id, &post.Title, &post.Body)
	return
}

func ShowPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))
	post, _ := GetPost(id)
	templates.ExecuteTemplate(w, "posts/show", post) 
}

func main() {
	router := httprouter.New()

	http.HandleFunc("/posts", PostIndex)
	http.HandleFunc("/posts/new", NewPost)
	http.HandleFunc("/posts/create", CreatePost)
	router.GET("/posts/:id", ShowPost)
	log.Fatal(http.ListenAndServe(":3000", router))
}