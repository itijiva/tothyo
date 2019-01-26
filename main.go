package main

import (
	"net/http"
	"html/template"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"strconv"
)

var templates = template.Must(template.New("t").ParseGlob("views/**/*.html"))

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

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	title := r.FormValue("title")
	body := r.FormValue("body")
	_, err := Db.Exec("update posts set title = $2, body = $3 where id = $1", id, title, body)
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

func ShowPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	post, _ := GetPost(id)
	templates.ExecuteTemplate(w, "posts/show", post) 
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	post, _ := GetPost(id)
	templates.ExecuteTemplate(w, "posts/edit", post) 
}

func RootPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.New("root.html").ParseFiles("views/root.html")

	tmpl.ExecuteTemplate(w, "root", "") 
}

func main() {
    r := mux.NewRouter()

    // This will serve files under http://localhost:8000/static/<filename>
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	
	r.HandleFunc("/", RootPage)
	r.HandleFunc("/posts", PostIndex).Methods("GET")
	r.HandleFunc("/posts/new", NewPost).Methods("GET")
	r.HandleFunc("/posts/create", CreatePost).Methods("POST")
	r.HandleFunc("/posts/{id:[0-9]+}", ShowPost).Methods("GET")
	r.HandleFunc("/posts/{id:[0-9]+}/edit", EditPost).Methods("GET")
	r.HandleFunc("/posts/{id:[0-9]+}/update", UpdatePost).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", r))
}