package main

type Post struct {
	Id int
	Title string
	Body string
}

func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("select id, title, body from posts where id = $1", id).Scan(&post.Id, &post.Title, &post.Body)
	return
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