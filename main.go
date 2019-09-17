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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
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
	// Routes consist of a path and a handler function.
	r.HandleFunc("/status", env.Status)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
