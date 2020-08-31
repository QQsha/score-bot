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

type FixtureUseCase struct {
	fixtureRepo repository.FixtureRepositoryInterface
	botRepo     repository.BotRepositoryInterface
	fixtureAPI  repository.APIRepositoryInterface
	langAPI     repository.DetectLanguageAPI
	log         *logrus.Logger
}

func NewFixtureUseCase(
	fixtureRepo repository.FixtureRepositoryInterface,
	botRepo repository.BotRepositoryInterface,
	fixtureAPI repository.APIRepositoryInterface,
	langAPI repository.DetectLanguageAPI,
) *FixtureUseCase {
	return &FixtureUseCase{
		fixtureRepo: fixtureRepo,
		botRepo:     botRepo,
		fixtureAPI:  fixtureAPI,
		langAPI:     langAPI,
		log:         logrus.New(),
	}
}

func (u FixtureUseCase) GetFixtures() {
	fixtures := u.fixtureAPI.GetFixtures()
	err := u.fixtureRepo.SaveFixtures(fixtures)
	if err != nil {
		u.log.Error(err)
	}
}
func (u FixtureUseCase) ZeroEventer(w http.ResponseWriter, r *http.Request) {
	fixture, err := u.fixtureRepo.NearestFixture(true)
	if err != nil {
		u.log.Error(err)
		return
	}
	u.fixtureRepo.ZeroEventer(fixture.FixtureID)
}
func (u FixtureUseCase) DateFix(w http.ResponseWriter, r *http.Request) {
	fixture := 571948
	u.fixtureRepo.DateFix(fixture)
}
func (u FixtureUseCase) TestPost(w http.ResponseWriter, r *http.Request) {
	fixture := u.fixtureAPI.GetFixtureDetails(571948)
	textMVP := u.GetWinners(fixture)
	// text := posts.CreateStatsPost(fixture)
	err := u.botRepo.SendPost(textMVP, nil)
	if err != nil {
		u.log.Error(err)
	}
}

func (u FixtureUseCase) GetLeaderboard() {
	winners, err := u.fixtureRepo.GetLeaderboard()
	if err != nil {
		u.log.Error(err)
		return
	}
	post := posts.CreateLeaderboardPost(winners)
	err = u.botRepo.SendPost(post, nil)
	if err != nil {
		u.log.Error(err)
	}
}

func (u FixtureUseCase) GetLineup(fixtureID, interval int) models.Lineup {
	lineup := models.Lineup{}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
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
			if len(fixtureDetail.API.Fixtures) == 0 {
				u.log.Errorf("fixture details is missing")
				return
			}
			if len(fixtureDetail.API.Fixtures[0].Events) > 0 {
				eventCount, err := u.fixtureRepo.EventCounter(fixtureID)
				if err != nil {
					u.log.Error(err)
					return
				}
				if len(fixtureDetail.API.Fixtures[0].Events) > eventCount {
					newEvents := fixtureDetail.API.Fixtures[0].Events[eventCount:]
					for _, event := range newEvents {
						if event.Type == "subst" && event.TeamID != 49 {
							u.fixtureRepo.EventIncrementer(fixtureID)
							continue
						}
						text := posts.CreatePostEvent(fixtureDetail, event)
						err := u.botRepo.SendPost(text, nil)
						if err != nil {
							u.log.Error(err)
						}
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
	text += u.GetWinners(fixtureDetail)
	err := u.botRepo.SendPost(text, nil)
	if err != nil {
		u.log.Error(err)
	}
	textStats := posts.CreateStatsPost(fixtureDetail)
	err = u.botRepo.SendPost(textStats, nil)
	if err != nil {
		u.log.Error(err)
	}
	textMVP, playersMVP := posts.CreateMVPPost(fixtureDetail)
	err = u.botRepo.SendPoll(textMVP, playersMVP)
	if err != nil {
		u.log.Error(err)
	}
}

func (u FixtureUseCase) GetRandomPrase(prases []string) string {
	post := prases[rand.Intn(len(prases))]
	return post
}

func (u FixtureUseCase) GetWinners(fixtureDetail models.FixtureDetails) string {
	ids, err := u.fixtureRepo.GetWinnersID(fixtureDetail)
	if err != nil {
		u.log.Error(err)
		return ""
	}
	winners := make([]models.User, 0)
	for _, id := range ids {
		winner, err := u.botRepo.GetChatUser(id)
		if err != nil {
			u.log.Error(err)
			return ""
		}
		err = u.fixtureRepo.AddLeader(winner, fixtureDetail.FixtureScore())
		if err != nil {
			u.log.Error(err)
			return ""
		}
		winners = append(winners, winner)
	}
	textWinners := posts.CreateWinnersPost(winners)
	return textWinners
}

func (u FixtureUseCase) LineupPoster() {
	u.GetFixtures()
	fixture, err := u.fixtureRepo.NearestFixture(true)
	if err != nil {
		u.log.Error(err)
		time.Sleep(30 * time.Minute)
		return
	}
	sleepTime := fixture.TimeTo - (53 * time.Minute)
	// sleepTime := time.Second * 3
	logrus.Info(fixture.HomeTeam.TeamName, " vs team: ", fixture.AwayTeam.TeamName, " id: ", fixture.FixtureID)
	logrus.Info("will send post after ", sleepTime)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(sleepTime, func() {
		lineup := u.GetLineup(fixture.FixtureID, 55)
		text := posts.CreatePostLineup(fixture, lineup)
		err := u.botRepo.SendPost(text, nil)
		if err != nil {
			u.log.Error(err)
		}
		err = u.fixtureRepo.FixturePosted(fixture.FixtureID)
		if err != nil {
			u.log.Error(err)
			wg.Done()
			return
		}
		wg.Done()
	})
	wg.Wait()
}
func (u FixtureUseCase) EventPoster() {
	fixture, err := u.fixtureRepo.NearestFixture(false)
	if err != nil {
		u.log.Error(err)
		time.Sleep(30 * time.Minute)
		return
	}
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

	fixture, err := u.fixtureRepo.NearestFixture(false)
	if err != nil {
		u.log.Error(err)
		return
	}
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
	//u.fixtureRepo.CreateTableLeagues()

	// u.fixtureRepo.CreateTableFixtures()
	// u.fixtureRepo.CreateTableLeaderboard()

}
func (u FixtureUseCase) CreateLeagues(w http.ResponseWriter, r *http.Request) {
	leagues := u.fixtureAPI.GetLeagues()
	err := u.fixtureRepo.AddLeague(leagues)
	if err != nil {
		u.log.Error(err)
		return
	}
}

func (u FixtureUseCase) MessageChecker() {
	chatID := u.botRepo.GetChatID()
	lastUpdateID, err := u.fixtureRepo.GetLastUpdate(chatID)
	if err != nil {
		u.log.Error(err)
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	sleepTime := time.Second * 3
	time.AfterFunc(sleepTime, func() {
		fmt.Println("start checking new messages")
		updates, err := u.botRepo.GetUpdates(lastUpdateID)
		if err != nil {
			u.log.Error(err)
			return
		}
		fmt.Printf("%+v\n", updates.Messages)

		spamWords, err := u.fixtureRepo.GetSpamWords()
		if err != nil {
			u.log.Error(err)
			return
		}
		englishPhrases := posts.GetEnglishPhrases()
		arsenalPhrases := posts.GetArsenalPhrases()

		var lastUpdate int
		for _, message := range updates.Messages {
			lastUpdate = message.UpdateID
			u.MessageHandler(message, spamWords, englishPhrases, arsenalPhrases)
		}
		if lastUpdate != 0 {
			err := u.fixtureRepo.MessageUpdates(chatID, lastUpdate)
			if err != nil {
				u.log.Error(err)
				return
			}
		}
		wg.Done()
	})
	wg.Wait()
}

func (u FixtureUseCase) MessageHandler(message models.FullMessage, spamWords []models.Spam, englishPhrases, arsPhrases []string) {
	if message.Message.Text == "" {
		message.Message.Text = message.Message.Caption
	}
	// commands check
	switch message.Message.Text {
	case "/rules@chelseaAntiSpamBot":
		rules := posts.GetRulePost()
		err := u.botRepo.SendPost(rules, nil)
		if err != nil {
			u.log.Error(err)
		}
		return
	case "/next@chelseaAntiSpamBot":
		fixture, err := u.fixtureRepo.NearestFixture(false)
		if err != nil {
			u.log.Error(err)
			return
		}
		post := posts.CreatePostNextGame(fixture)
		err = u.botRepo.SendPost(post, nil)
		if err != nil {
			u.log.Error(err)
		}
		return
	case "/leaderboard@chelseaAntiSpamBot":
		u.GetLeaderboard()
		return

	case "/predict@chelseaAntiSpamBot":
		err := u.botRepo.DeleteMessage(message.Message.MessageID)
		if err != nil {
			u.log.Error(err)
		}
		return
	}
	// predict check
	if strings.Contains(strings.ToLower(message.Message.Text), "/predict") {
		var homeTeamScore, awayTeamScore int
		score := strings.Split(message.Message.Text, " ")
		errorMsg := "wrong text format, use 3-1, where 3 is home team, and 1 is away team"
		if len(score) != 2 {
			err := u.botRepo.SendPost(errorMsg, &message.Message.MessageID)
			if err != nil {
				u.log.Error(err)
			}
			return
		}
		results := strings.Split(score[1], "-")
		if len(results) != 2 {
			err := u.botRepo.SendPost(errorMsg, &message.Message.MessageID)
			if err != nil {
				u.log.Error(err)
			}
			return
		}
		homeTeamScore, err := strconv.Atoi(results[0])
		if err != nil {
			err := u.botRepo.SendPost(errorMsg, &message.Message.MessageID)
			if err != nil {
				u.log.Error(err)
			}
			return
		}
		awayTeamScore, err = strconv.Atoi(results[1])
		if err != nil {
			err := u.botRepo.SendPost(errorMsg, &message.Message.MessageID)
			if err != nil {
				u.log.Error(err)
			}
			return
		}
		fixture, err := u.fixtureRepo.NearestFixture(false)
		if err != nil {
			u.log.Error(err)
			return
		}
		err = u.fixtureRepo.AddPrediction(fixture.FixtureID, message.Message.From.ID, homeTeamScore, awayTeamScore)
		if err != nil {
			u.log.Error(err)
			return
		}
		count, err := u.fixtureRepo.PredictionCounter(fixture.FixtureID)
		if err != nil {
			u.log.Error(err)
			return
		}
		postText := "Your predict, " + fixture.HomeTeam.TeamName + " " + strconv.Itoa(homeTeamScore) + " - " + strconv.Itoa(awayTeamScore) + " " + fixture.AwayTeam.TeamName + " is accepted!\n"
		postText += "Total predictions on this match: " + strconv.Itoa(count)
		err = u.botRepo.SendPost(postText, &message.Message.MessageID)
		if err != nil {
			u.log.Error(err)
		}
		return
	}
	// spam check
	isSpam, banDuration := u.IsSpam(message, spamWords)
	if isSpam {
		err := u.botRepo.RestrictUser(message.Message.From.ID, banDuration)
		if err != nil {
			u.log.Error(err)
		}
		err = u.botRepo.DeleteMessage(message.Message.MessageID)
		if err != nil {
			u.log.Error(err)
		}
		return
	}
	// emoji check
	if u.IsEmoji(message) {
		err := u.botRepo.DeleteMessage(message.Message.MessageID)
		if err != nil {
			u.log.Error(err)
		}
	}
	// arsenal check
	if strings.Contains(strings.ToLower(message.Message.Text), "arsenal") {
		randNum := rand.Intn(3)
		if randNum%3 != 0 {
			return
		}
		text := u.GetRandomPrase(arsPhrases)
		err := u.botRepo.SendPost(text, &message.Message.MessageID)
		if err != nil {
			u.log.Error(err)
		}
		return
	}
	// check msg on names
	withNames := u.langAPI.NameDetector(message.Message.Text)
	if !withNames {
		// english check
		notEnglish, _ := u.langAPI.EnglisheDetector(message.Message.Text)
		if notEnglish {
			text := u.GetRandomPrase(englishPhrases)
			err := u.botRepo.SendPost(text, &message.Message.MessageID)
			if err != nil {
				u.log.Error(err)
			}
		}
	}
	// tx := u.langAPI.EnglishDetectorTest(message.Message.Text)
	// tx := u.langAPI.NameDetectorTest(message.Message.Text)
	// u.botRepo.SendPost(tx, &message.Message.MessageID)
}

func (u FixtureUseCase) IsSpam(text models.FullMessage, spamWords []models.Spam) (bool, int) {
	if text.Message.ForwardFromChat.ID == -1001044276483 || text.Message.From.Username == "qqshaa" || text.Message.From.Username == "KingSuperFrank" {
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
func (u FixtureUseCase) IsEmoji(text models.FullMessage) bool {
	if text.Message.Dice.Emoji != "" {
		return true
	}
	return false
}

func (u FixtureUseCase) GetSpamWordsHandler(w http.ResponseWriter, r *http.Request) {
	spamWords, err := u.fixtureRepo.GetSpamWords()
	if err != nil {
		u.log.Error(err)
		return
	}
	res, _ := json.Marshal(spamWords)
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, string(res))
}

func (u FixtureUseCase) AddNewSpamWordHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10)
	if err != nil {
		u.log.Error(err)
		return
	}
	ban := r.FormValue("ban")
	spam := r.FormValue("spam")
	banDur, err := strconv.Atoi(ban)
	if err != nil {
		u.log.Error(err)
		return
	}
	err = u.fixtureRepo.AddSpamWord(spam, banDur)
	if err != nil {
		u.log.Error(err)
		return
	}
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
	err := u.fixtureRepo.DeleteSpamWord(spamWord)
	if err != nil {
		u.log.Error(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, "")
}
func (u FixtureUseCase) ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./react-api/build/index.html")
}
