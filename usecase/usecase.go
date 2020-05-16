package usecase

import (
	// "fmt"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QQsha/score-bot/models"
	"github.com/QQsha/score-bot/repository"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "15:04"
)

type FixtureUseCase struct {
	fixtureRepo repository.FixtureRepository
	botRepo     repository.BotRepository
	fixtureAPI  repository.APIRepositoryInterface
}

func NewFixtureUseCase(
	fixtureRepo repository.FixtureRepository,
	botRepo repository.BotRepository,
	fixtureAPI repository.APIRepositoryInterface,
) *FixtureUseCase {
	return &FixtureUseCase{
		fixtureRepo: fixtureRepo,
		botRepo:     botRepo,
		fixtureAPI:  fixtureAPI,
	}
}

func (u FixtureUseCase) GetFixtures() {
	fixtures := u.fixtureAPI.GetFixtures()
	u.fixtureRepo.SaveFixtures(fixtures)
}

func (u FixtureUseCase) GetLineup(fixtureID int) models.Lineup {
	lineup := models.Lineup{}
	ticker := time.NewTicker(10 * time.Second)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				lineup = u.fixtureAPI.GetLineup(fixtureID)
			}
			if len(lineup.API.LineUps.Chelsea.StartXI) != 0 {
				stop <- true
			}
		}
	}()
	<-stop
	ticker.Stop()
	return lineup
}
func (u FixtureUseCase) NewEventChecker(fixtureID int) {
	ticker := time.NewTicker(time.Second * 3)
	fixtureDetail := models.FixtureDetails{}
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				fixtureDetail = u.fixtureAPI.GetFixtureDetails(fixtureID)
				fmt.Println(fixtureDetail.API.Fixtures[0].StatusShort)
				if len(fixtureDetail.API.Fixtures) > 0 && len(fixtureDetail.API.Fixtures[0].Events) > 0 {
					eventCount := u.fixtureRepo.EventCounter(fixtureID)
					if len(fixtureDetail.API.Fixtures[0].Events) > eventCount {
						newEvents := fixtureDetail.API.Fixtures[0].Events[eventCount:]
						for _, event := range newEvents {
							text := u.CreatePostEvent(event, fixtureDetail.API.Fixtures[0].Score.Fulltime)
							u.botRepo.SendPost(text)
							u.fixtureRepo.EventIncrementer(fixtureID)
						}
					}
				}
				if fixtureDetail.API.Fixtures[0].StatusShort == "FT" {
					stop <- true
					return
				}
			}
		}
	}()
	<-stop
	ticker.Stop()
	text := u.CreateFullTimePost(fixtureDetail)
	u.botRepo.SendPost(text)
	return
}

func (u FixtureUseCase) CreatePostLineup(fixture models.Fixture, lineup models.Lineup) string {
	text := "📣💙*Match of the day:*💙📣\n"
	text += "*" + fixture.HomeTeam.TeamName + " - " + fixture.AwayTeam.TeamName + "*\n"
	text += "\n *Line-up (" + lineup.API.LineUps.Chelsea.Formation + "):*\n"
	for _, player := range lineup.API.LineUps.Chelsea.StartXI {
		text += player.Player + " (" + strconv.Itoa(player.Number) + ")" + "\n"
	}
	text += "\n *Substitutes*:\n"
	for _, player := range lineup.API.LineUps.Chelsea.Substitutes {
		text += player.Player + " (" + strconv.Itoa(player.Number) + ")" + "\n"
	}
	text += "\n *Match starts today at*:\n"
	loc, _ := time.LoadLocation("Asia/Tehran")
	text += "🇮🇷*Tehran*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Africa/Lagos")
	text += "🇳🇬*Abuja*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Europe/Moscow")
	text += "🇷🇺*Moscow*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Asia/Almaty")
	text += "🇰🇿*Nur-Sultan*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Europe/London")
	text += "🇬🇧*London*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	return text
}
func (u FixtureUseCase) CreatePostEvent(event models.Event, currentScore string) string {
	var post string
	switch event.Type {
	// "elapsed": 46,
	// "elapsed_plus": null,
	// "team_id": 49,
	// "teamName": "Chelsea",
	// "player_id": 2285,
	// "player": "A. Rüdiger",
	// "assist_id": 19220,
	// "assist": "M. Mount",
	// "type": "Goal",
	// "detail": "Normal Goal",
	// "comments": null
	case "Goal":
		post = "Goal by " + event.Player + ", assist by " + event.Assist + " current score: " + currentScore
	case "Card":
		post = "someone got card"
	case "subst":
		post = "substa was made"
	}
	return post
}
func (u FixtureUseCase) CreateFullTimePost(fixture models.FixtureDetails) string {
	var post string
	post = "full time score: " + fixture.API.Fixtures[0].Score.Fulltime
	return post
}
func (u FixtureUseCase) LineupPoster() {
	u.GetFixtures()
	fixture := u.fixtureRepo.NearestFixture()
	sleepTime := fixture.TimeTo - (55 * time.Minute)
	// sleepTime := time.Second * 3
	logrus.Info(fixture.HomeTeam.TeamName, " vs team: ", fixture.AwayTeam.TeamName, " id: ", fixture.FixtureID)
	logrus.Info("will send post after ", sleepTime)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(sleepTime, func() {
		lineup := u.GetLineup(fixture.FixtureID)
		text := u.CreatePostLineup(fixture, lineup)
		u.botRepo.SendPost(text)
		// u.fixtureRepo.FixturePosted(fixture.FixtureID)
		wg.Done()
	})
	wg.Wait()
}
func (u FixtureUseCase) EventPoster() {
	fixture := u.fixtureRepo.NearestFixture()
	// sleepTime := fixture.TimeTo + (5 * time.Minute)
	sleepTime := time.Hour * 3
	wg := &sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(sleepTime, func() {
		fmt.Println("start checking new events")
		u.NewEventChecker(fixture.FixtureID)
		wg.Done()
	})
	wg.Wait()
}

func (u FixtureUseCase) Status(w http.ResponseWriter, r *http.Request) {
	status := u.fixtureAPI.StatusCheck()
	reqLeft := strconv.Itoa(status.API.Status.RequestsLeft)

	fixture := u.fixtureRepo.NearestFixture()
	io.WriteString(w, "Next match: "+fixture.HomeTeam.TeamName+" vs "+fixture.AwayTeam.TeamName+"\n")
	timeLeft := (fixture.TimeTo - (time.Minute * 55)).String()
	io.WriteString(w, "Post will be in: "+timeLeft+"\n")
	io.WriteString(w, "Fixture ID: "+strconv.Itoa(fixture.FixtureID)+"\n")
	logrus.Info("vs team: ", fixture.HomeTeam.TeamName, fixture.AwayTeam.TeamName, " id: ", fixture.FixtureID)
	logrus.Info("will send post after ", fixture.TimeTo-(time.Minute*55))
	io.WriteString(w, "Requests remaining:"+reqLeft)
}

func (u FixtureUseCase) CreateTable(w http.ResponseWriter, r *http.Request) {
	u.fixtureRepo.CreateTableSpamBase()
	u.fixtureRepo.CreateTableLastUpdate()
	u.fixtureRepo.CreateTableFixtures()

}

func (u FixtureUseCase) MessageChecker() {
	var offset int
	chatID := u.botRepo.GetChatID()
	offset = u.fixtureRepo.GetLastUpdate(chatID)
	sleepTime := time.Second * 5
	wg := &sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(sleepTime, func() {
		fmt.Println("start checking new messages")
		updates := u.botRepo.GetUpdates(offset)
		fmt.Printf("%+v\n", updates.Messages)
		spamWords := u.fixtureRepo.GetSpamWords()
		var lastUpdate int
		for _, message := range updates.Messages {
			isSpam, banDuration := IsSpam(message.Message.Text, spamWords)
			if isSpam {
				fmt.Println("SPAM")
				u.botRepo.RestrictUser(message.Message.From.ID, banDuration)
			}
			lastUpdate = message.UpdateID
		}
		if lastUpdate != 0 {
			u.fixtureRepo.MessageUpdates(chatID, lastUpdate)
		}
		wg.Done()
	})
	wg.Wait()
}
func IsSpam(text string, spamWords []models.Spam) (bool, int) {
	for _, word := range spamWords {
		if strings.Contains(text, word.Word) {
			return true, word.BanDuration
		}
	}
	return false, 0
}

func (u FixtureUseCase) GetSpamWordsHandler(w http.ResponseWriter, r *http.Request) {
	spamWords := u.fixtureRepo.GetSpamWords()
	res, _ := json.Marshal(spamWords)
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, string(res))
}

func (u FixtureUseCase) AddNewSpamWordHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10)
	ban := r.FormValue("ban")
	spam := r.FormValue("spam")
	banDur, err := strconv.Atoi(ban)
	if err != nil {
		fmt.Println(err)
	}
	u.fixtureRepo.AddSpamWord(spam, banDur)
	// res, _ := json.Marshal(spamWords)
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, "")
}

func (u FixtureUseCase) TestHandler(w http.ResponseWriter, r *http.Request) {
	body := []byte{}
	r.Body.Read(body)
	fmt.Println(string(body))
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, string("TEST HANDLER"))
}

func (u FixtureUseCase) DeleteSpamWordHandler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	spamWord := keys.Get("spam")

	fmt.Println(spamWord)
	u.fixtureRepo.DeleteSpamWord(spamWord)
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, "")
}
func (u FixtureUseCase) ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./react-api/build/index.html")
}
