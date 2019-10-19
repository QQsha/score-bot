package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"score_bot/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Env struct {
	DB *sql.DB
}

const (
	host   = "0.0.0.0"
	port   = 5432
	user   = "postgres"
	dbname = "qqsha"
)

func main() {
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
	// 	host, port, user, dbname)
	// db, err := sql.Open("postgres", psqlInfo)
	db, err := sql.Open("postgres", "postgres://xsqidgwwvwvgkm:1e82bd5c5b23996ee1ed11dfaa89447adc5c524999c574b6c24b67c0c1a22604@ec2-75-101-153-56.compute-1.amazonaws.com:5432/ddgva2m0b3akm5")
	if err != nil {
		panic(err)
	}
	env := &handlers.Env{DB: db}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	go env.SetUp()

	r := mux.NewRouter()
	// Routes consist of a path and a handler function  .
	r.HandleFunc("/", env.Status)
	r.HandleFunc("/create", env.CreateTable)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(handlers.DetermineListenAddress(), r))
}
