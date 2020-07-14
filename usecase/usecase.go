package usecase

import (
	// "fmt"

	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QQsha/score-bot/models"
	"github.com/QQsha/score-bot/posts"
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
	langAPI     repository.DetectLanguageAPI
}

func NewFixtureUseCase(
	fixtureRepo repository.FixtureRepository,
	botRepo repository.BotRepository,
	fixtureAPI repository.APIRepositoryInterface,
	langAPI repository.DetectLanguageAPI,
) *FixtureUseCase {
	return &FixtureUseCase{
		fixtureRepo: fixtureRepo,
		botRepo:     botRepo,
		fixtureAPI:  fixtureAPI,
		langAPI:     langAPI,
	}
}

func (u FixtureUseCase) GetFixtures() {
	fixtures := u.fixtureAPI.GetFixtures()
	u.fixtureRepo.SaveFixtures(fixtures)
}
func (u FixtureUseCase) ZeroEventer(w http.ResponseWriter, r *http.Request) {
	fixture := u.fixtureRepo.NearestFixture(true)
	u.fixtureRepo.ZeroEventer(fixture.FixtureID)
}
func (u FixtureUseCase) TestPost(w http.ResponseWriter, r *http.Request) {
	fixture := u.fixtureAPI.GetFixtureDetails(333)
	textMVP, playersMVP := posts.CreateMVPPost(fixture)
	u.botRepo.SendPoll(textMVP, playersMVP)
	// text := posts.CreateStatsPost(fixture)
	// u.botRepo.SendPost(text, nil)
}
func (u FixtureUseCase) GetLineup(fixtureID int) models.Lineup {
	lineup := models.Lineup{}
	ticker := time.NewTicker(40 * time.Second)
	stop := make(chan bool)
	go func() {
		for {
			<-ticker.C
			lineup = u.fixtureAPI.GetLineup(fixtureID)
			if len(lineup.API.LineUps.Chelsea.StartXI) == 11 {
				ticker.Stop()
				stop <- true
			}
		}
	}()
	<-stop
	return lineup
}
func (u FixtureUseCase) NewEventChecker(fixtureID int) {
	ticker := time.NewTicker(time.Minute * 2)
	fixtureDetail := models.FixtureDetails{}
	stop := make(chan bool)
	go func() {
		for {
			<-ticker.C
			fixtureDetail = u.fixtureAPI.GetFixtureDetails(fixtureID)
			if len(fixtureDetail.API.Fixtures) > 0 && len(fixtureDetail.API.Fixtures[0].Events) > 0 {
				eventCount := u.fixtureRepo.EventCounter(fixtureID)
				if len(fixtureDetail.API.Fixtures[0].Events) > eventCount {
					newEvents := fixtureDetail.API.Fixtures[0].Events[eventCount:]
					for _, event := range newEvents {
						if event.Type == "subst" && event.TeamID != 49 {
							u.fixtureRepo.EventIncrementer(fixtureID)
							continue
						}
						text := posts.CreatePostEvent(fixtureDetail, event)
						u.botRepo.SendPost(text, nil)
						u.fixtureRepo.EventIncrementer(fixtureID)
					}
				}
				// if goal was canceled by var
				if eventCount > len(fixtureDetail.API.Fixtures[0].Events) {
					u.fixtureRepo.EventDecrementer(fixtureID)
				}
			}
			if fixtureDetail.API.Fixtures[0].StatusShort == "FT" {
				ticker.Stop()
				stop <- true
				return
			}
		}
	}()
	<-stop
	text := posts.CreateFullTimePost(fixtureDetail)
	u.botRepo.SendPost(text, nil)
	textStats := posts.CreateStatsPost(fixtureDetail)
	u.botRepo.SendPost(textStats, nil)
	textMVP, playersMVP := posts.CreateMVPPost(fixtureDetail)
	u.botRepo.SendPoll(textMVP, playersMVP)
}

func (u FixtureUseCase) GetRandomPrase(prases []string) string {
	post := prases[rand.Intn(len(prases))]
	return post
}

func (u FixtureUseCase) LineupPoster() {
	u.GetFixtures()
	fixture := u.fixtureRepo.NearestFixture(true)
	sleepTime := fixture.TimeTo - (55 * time.Minute)
	// sleepTime := time.Second * 3
	logrus.Info(fixture.HomeTeam.TeamName, " vs team: ", fixture.AwayTeam.TeamName, " id: ", fixture.FixtureID)
	logrus.Info("will send post after ", sleepTime)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(sleepTime, func() {
		lineup := u.GetLineup(fixture.FixtureID)
		text := posts.CreatePostLineup(fixture, lineup)
		u.botRepo.SendPost(text, nil)
		u.fixtureRepo.FixturePosted(fixture.FixtureID)
		wg.Done()
	})
	wg.Wait()
}
func (u FixtureUseCase) EventPoster() {
	fixture := u.fixtureRepo.NearestFixture(false)
	sleepTime := fixture.TimeTo + (5 * time.Minute)
	// sleepTime := time.Second * 2
	wg := &sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(sleepTime, func() {
		fmt.Println("start checking new Events")
		u.NewEventChecker(fixture.FixtureID)
		wg.Done()
	})
	wg.Wait()
}

func (u FixtureUseCase) Status(w http.ResponseWriter, r *http.Request) {
	status := u.fixtureAPI.StatusCheck()
	reqLeft := strconv.Itoa(status.API.Status.RequestsLeft)

	fixture := u.fixtureRepo.NearestFixture(false)
	io.WriteString(w, "Next match: "+fixture.HomeTeam.TeamName+" vs "+fixture.AwayTeam.TeamName+"\n")
	timeLeft := (fixture.TimeTo - (time.Minute * 55)).String()
	io.WriteString(w, "Post will be in: "+timeLeft+"\n")
	io.WriteString(w, "Fixture ID: "+strconv.Itoa(fixture.FixtureID)+"\n")
	logrus.Info("vs team: ", fixture.HomeTeam.TeamName, fixture.AwayTeam.TeamName, " id: ", fixture.FixtureID)
	logrus.Info("will send post after ", fixture.TimeTo-(time.Minute*55))
	io.WriteString(w, "Requests remaining:"+reqLeft)
}

func (u FixtureUseCase) CreateTable(w http.ResponseWriter, r *http.Request) {
	// u.fixtureRepo.CreateTableSpamBase()
	// u.fixtureRepo.CreateTableLastUpdate()
	// u.fixtureRepo.CreateTableFixtures()
	u.fixtureRepo.CreateTableLeagues()

}
func (u FixtureUseCase) CreateLeagues(w http.ResponseWriter, r *http.Request) {
	leagues := u.fixtureAPI.GetLeagues()
	u.fixtureRepo.AddLeague(leagues)
}

func (u FixtureUseCase) MessageChecker() {
	chatID := u.botRepo.GetChatID()
	lastUpdateID := u.fixtureRepo.GetLastUpdate(chatID)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	sleepTime := time.Second * 3
	time.AfterFunc(sleepTime, func() {
		fmt.Println("start checking new messages")
		updates := u.botRepo.GetUpdates(lastUpdateID)
		fmt.Printf("%+v\n", updates.Messages)

		spamWords := u.fixtureRepo.GetSpamWords()
		englishPhrases := posts.GetEnglishPhrases()
		arsenalPhrases := posts.GetArsenalPhrases()

		var lastUpdate int
		for _, message := range updates.Messages {
			lastUpdate = message.UpdateID
			u.MessageHandler(message, spamWords, englishPhrases, arsenalPhrases)
		}
		if lastUpdate != 0 {
			u.fixtureRepo.MessageUpdates(chatID, lastUpdate)
		}
		wg.Done()
	})
	wg.Wait()
}

func (u FixtureUseCase) MessageHandler(message models.Message, spamWords []models.Spam, englishPhrases, arsPhrases []string) {
	if message.Message.Text == "" {
		message.Message.Text = message.Message.Caption
	}
	// commands check
	switch message.Message.Text {
	case "/rules@chelseaAntiSpamBot":
		rules := posts.GetRulePost()
		u.botRepo.SendPost(rules, nil)
		return
	case "/next@chelseaAntiSpamBot":
		fixture := u.fixtureRepo.NearestFixture(false)
		post := posts.CreatePostNextGame(fixture)
		u.botRepo.SendPost(post, nil)
		return
	}
	// spam check
	isSpam, banDuration := IsSpam(message, spamWords)
	if isSpam {
		u.botRepo.RestrictUser(message.Message.From.ID, banDuration)
		u.botRepo.DeleteMessage(message.Message.MessageID)
		fmt.Println("SPAM")
		return
	}
	// arsenal check
	if strings.Contains(strings.ToLower(message.Message.Text), "arsenal") {
		text := u.GetRandomPrase(arsPhrases)
		u.botRepo.SendPost(text, &message.Message.MessageID)
		return
	}
	// check msg on names
	withNames := u.langAPI.NameDetector(message.Message.Text)
	if !withNames {
		// english check
		notEnglish, _ := u.langAPI.EnglisheDetector(message.Message.Text)
		if notEnglish {
			text := u.GetRandomPrase(englishPhrases)
			u.botRepo.SendPost(text, &message.Message.MessageID)
		}
	}
	// tx := u.langAPI.EnglishDetectorTest(message.Message.Text)
	// tx := u.langAPI.NameDetectorTest(message.Message.Text)
	// u.botRepo.SendPost(tx, &message.Message.MessageID)
}

func IsSpam(text models.Message, spamWords []models.Spam) (bool, int) {
	if text.Message.ForwardFromChat.ID == -1001044276483 || text.Message.From.Username == "qqshaaa" {
		return false, 0
	}
	for _, ent := range text.Message.CaptionEntities {
		if ent.Type == "mention" {
			return true, 1
		}
	}
	if text.Message.ForwardFrom.ID != 0 || text.Message.ForwardFromChat.ID != 0 {
		return true, 1
	}
	for _, word := range spamWords {
		if strings.Contains(strings.ToLower(text.Message.Text), word.Word) {
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
func (u FixtureUseCase) MashaAanswer(w http.ResponseWriter, r *http.Request) {
	text := `
	спасиб0 биб буп 
	`
	ss := 83
	u.botRepo.SendPost(text, &ss)
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, "")
}
