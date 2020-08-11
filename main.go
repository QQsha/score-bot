package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/QQsha/score-bot/repository"
	"github.com/QQsha/score-bot/usecase"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	
)

func main() {
	postgresConn := os.Getenv("POSTGRES_SCOREBOT")

	db, teardown, err := repository.Connect(postgresConn)
	if err != nil {
		log.Fatal(err)
	}
	defer teardown()
	channelID := "@Chelsea"
	// channelID := "@Chelseafuns"
	fixtureBotToken := os.Getenv("SCORE_BOT")

	chatID := strconv.Itoa(-1001490460294) //"Chelsea chat"
	// chatID := strconv.Itoa(-1001276457176) //"Группа а не канал"
	antiSpamBotToken := os.Getenv("ANTISPAM_BOT")
	antiSpamBotRepo := repository.NewBotRepository(antiSpamBotToken, chatID)
	fixtureBotRepo := repository.NewBotRepository(fixtureBotToken, channelID)
	fixtureRepo := repository.NewFixtureRepository(db)
	detectLangApi := repository.NewLanguageAPI()
	apiToken := os.Getenv("APIFOOTBALL_TOKEN")
	apiRepo := repository.NewAPIRepository(apiToken)
	// apiRepoMOCK := mock.NewAPIRepositoryMock()

	fixtureBot := usecase.NewFixtureUseCase(*fixtureRepo, *fixtureBotRepo, *apiRepo, *detectLangApi)
	antiSpamBot := usecase.NewFixtureUseCase(*fixtureRepo, *antiSpamBotRepo, *apiRepo, *detectLangApi)

	go func() {
		for {
			fixtureBot.LineupPoster()
		}
	}()
	go func() {
		for {
			fixtureBot.EventPoster()
		}
	}()
	go func() {
		for {
			antiSpamBot.MessageChecker()
		}
	}()
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./react-api/build/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/stats", fixtureBot.Status)
	r.HandleFunc("/create_table", fixtureBot.CreateTable)
	r.HandleFunc("/create_leagues", fixtureBot.CreateLeagues)
	r.HandleFunc("/get_spam", antiSpamBot.GetSpamWordsHandler)
	r.HandleFunc("/add_new_spam", antiSpamBot.AddNewSpamWordHandler)
	r.HandleFunc("/delete_spam", antiSpamBot.DeleteSpamWordHandler)
	r.HandleFunc("/zero_event", fixtureBot.ZeroEventer)
	r.HandleFunc("/date_fix", fixtureBot.DateFix)
	r.HandleFunc("/test", fixtureBot.TestPost)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./react-api/build/index.html")
	})
	handler := cors.Default().Handler(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
