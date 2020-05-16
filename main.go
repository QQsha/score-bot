package main

import (
	"log"
	"net/http"
	"os"

	"github.com/QQsha/score-bot/repository"
	mock "github.com/QQsha/score-bot/repository/mock"
	"github.com/QQsha/score-bot/usecase"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

const (
	// host   = "0.0.0.0"
	// port   = 5432
	user   = "postgres"
	dbname = "qqsha"
)

func main() {
	postgresConn := os.Getenv("POSTGRES_SCOREBOT")
	// portDB := os.Getenv("DB_PORT")
	// if portDB == "" {
	// 	portDB = "5432"
	// }

	// host := os.Getenv("DB_HOST")
	// if host == "" {
	// 	host = "0.0.0.0"
	// }
	// postgresConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
	// 	host, portDB, user, dbname)
	db, teardown, err := repository.Connect(postgresConn)
	if err != nil {
		log.Fatal(err)
	}
	defer teardown()
	logrus.Info("Successfully connected!")
	// chatID := "@Chelsea"
	testChannelID := "@Chelseafuns"
	fixtureBotToken := os.Getenv("SCORE_BOT")
	testGroupID := "Группа а не канал"
	antiSpamBotToken := "1238653037:AAE5szQlripWfNcHPVnBF2UasrDlZ-KS_nk"
	antiSpamBotRepo := repository.NewBotRepository(antiSpamBotToken, testGroupID)
	fixtureBotRepo := repository.NewBotRepository(fixtureBotToken, testChannelID)
	fixtureRepo := repository.NewFixtureRepository(db)
	// apiToken := "9128771ca86462be53b41393a002341e"
	// apiRepo := repository.NewAPIRepository(apiToken)
	apiRepo := mock.NewAPIRepositoryMock()
	fixtureBot := usecase.NewFixtureUseCase(*fixtureRepo, *fixtureBotRepo, *apiRepo)
	antiSpamBot := usecase.NewFixtureUseCase(*fixtureRepo, *antiSpamBotRepo, *apiRepo)
	// go func() {
	// 	for {
	// 		fixtureBot.LineupPoster()
	// 	}
	// }()
	// go func() {
	// 	for {
	// 		fixtureBot.EventPoster()
	// 	}
	// }()
	go func() {
		for {
			antiSpamBot.MessageChecker()
		}
	}()
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./react-api/build/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/", fixtureBot.Status)
	r.HandleFunc("/create_table", fixtureBot.CreateTable)
	r.HandleFunc("/get_spam", antiSpamBot.GetSpamWordsHandler)
	r.HandleFunc("/add_new_spam", antiSpamBot.AddNewSpamWordHandler)
	r.HandleFunc("/delete_spam", antiSpamBot.DeleteSpamWordHandler)
	r.HandleFunc("/spam", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./react-api/build/index.html")
	})
	handler := cors.Default().Handler(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
