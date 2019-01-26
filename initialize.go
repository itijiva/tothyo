package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=postgres dbname=gwp password=gwp sslmode=disable")
	
	if err != nil {
		panic(err)
	}
}