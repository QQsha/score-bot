package main

import (
	"log"
	"net/http"
	"os"

	"database/sql"

	"score_bot/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Env struct {
	DB *sql.DB
}

// const (
// 	host   = "0.0.0.0"
// 	port   = 5432
// 	user   = "postgres"
// 	dbname = "qqsha"
// )

func main() {
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
	// 	host, port, user, dbname)
	// db, err := sql.Open("postgres", psqlInfo)
	
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_SCOREBOT"))
	if err != nil {
		panic(err)
	}
	env := &handlers.Env{DB: db}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	logrus.Info("Successfully connected!")
	go env.Start()
	go handlers.Updater()

	r := mux.NewRouter()
	// Routes consist of a path and a handler function  .
	r.HandleFunc("/", env.Status)
	r.HandleFunc("/create", env.CreateTable)
	// Bind to a port and pass our router in
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
