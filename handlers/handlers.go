package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	timeFormat   = "15:04 2006-01-02"
	keyHeader    = "X-RapidAPI-Key"
	tokenAPIFoot = "9128771ca86462be53b41393a002341e"
	// ddbURI       = "postgres://xsqidgwwvwvgkm:1e82bd5c5b23996ee1ed11dfaa89447adc5c524999c574b6c24b67c0c1a22604@ec2-75-101-153-56.compute-1.amazonaws.com:5432/ddgva2m0b3akm5"
	botToken = "840859313:AAFfNUxxiaw6MIj9_5XSIeelJv7gns8qRqk"
	oldBot   = "569665229:AAFFOoITLtgjpxsWtAoHTATMNv5mex53JXU"
	chatID   = "@Chelsea"
	testChat = "@chelseafuns"
)

// https://server1.api-football.com/fixtures/team/49
type Fixtures struct {
	API struct {
		Results  int       `json:"results"`
		Fixtures []Fixture `json:"fixtures"`
	} `json:"api"`
}
type Fixture struct {
	FixtureID       int         `json:"fixture_id"`
	LeagueID        int         `json:"league_id"`
	EventDate       time.Time   `json:"event_date"`
	EventTimestamp  int         `json:"event_timestamp"`
	FirstHalfStart  interface{} `json:"firstHalfStart"`
	SecondHalfStart interface{} `json:"secondHalfStart"`
	Round           string      `json:"round"`
	Status          string      `json:"status"`
	StatusShort     string      `json:"statusShort"`
	Elapsed         int         `json:"elapsed"`
	Venue           string      `json:"venue"`
	Referee         interface{} `json:"referee"`
	HomeTeam        struct {
		TeamID   int    `json:"team_id"`
		TeamName string `json:"team_name"`
		Logo     string `json:"logo"`
	} `json:"homeTeam"`
	AwayTeam struct {
		TeamID   int    `json:"team_id"`
		TeamName string `json:"team_name"`
		Logo     string `json:"logo"`
	} `json:"awayTeam"`
	GoalsHomeTeam int `json:"goalsHomeTeam"`
	GoalsAwayTeam int `json:"goalsAwayTeam"`
	Score         struct {
		Halftime  string      `json:"halftime"`
		Fulltime  string      `json:"fulltime"`
		Extratime interface{} `json:"extratime"`
		Penalty   interface{} `json:"penalty"`
	} `json:"score"`
	TimeTo time.Duration
}
type Lineup struct {
	API struct {
		Results int `json:"results"`
		LineUps struct {
			Chelsea struct {
				Formation string `json:"formation"`
				StartXI   []struct {
					TeamID   int    `json:"team_id"`
					PlayerID int    `json:"player_id"`
					Player   string `json:"player"`
					Number   int    `json:"number"`
					Pos      string `json:"pos"`
				} `json:"startXI"`
				Substitutes []struct {
					TeamID   int    `json:"team_id"`
					PlayerID int    `json:"player_id"`
					Player   string `json:"player"`
					Number   int    `json:"number"`
					Pos      string `json:"pos"`
				} `json:"substitutes"`
				Coach string `json:"coach"`
			} `json:"Chelsea"`
		} `json:"lineUps,omitempty"`
	} `json:"api"`
}
type Status struct {
	API struct {
		Code    int `json:"code"`
		Results int `json:"results"`
		Status  struct {
			User             string        `json:"user"`
			Email            string        `json:"email"`
			Plan             string        `json:"plan"`
			Token            string        `json:"token"`
			Active           string        `json:"active"`
			SubscriptionEnd  time.Time     `json:"subscription_end"`
			Requests         int           `json:"requests"`
			RequestsLimitDay int           `json:"requests_limit_day"`
			Payments         []interface{} `json:"payments"`
			RequestsLeft     int
		} `json:"status"`
	} `json:"api"`
}
type Env struct {
	DB *sql.DB
}

func (env *Env) GetFixtures(fixtureID int) {
	now := time.Now()
	status := env.StatusCheck()
	if status.API.Status.RequestsLeft == 0 {

		logrus.Info("reached limit request")
		time.Sleep(time.Hour)
		env.GetFixtures(fixtureID)
	}
	uri := "https://server1.api-football.com/fixtures/team/49"

	client := http.Client{}

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}
	headers := http.Header{
		keyHeader: []string{tokenAPIFoot},
	}
	request.Header = headers
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fixtures := Fixtures{}
	err = json.Unmarshal(body, &fixtures)
	if err != nil {
		log.Fatalln(err)
	}
	if len(fixtures.API.Fixtures) > 0 {
		for _, fixture := range fixtures.API.Fixtures {
			if fixture.StatusShort == "NS" && fixture.FixtureID != fixtureID {
				query := `
					INSERT INTO fixtures (id, league_id, date, home_team, away_team)
					VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT (id)
					DO
					UPDATE
					SET date = EXCLUDED.date`
				_, err = env.DB.Exec(
					query, fixture.FixtureID, fixture.LeagueID, fixture.EventDate, fixture.HomeTeam.TeamName, fixture.AwayTeam.TeamName)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	query := `
		DELETE FROM fixtures
		WHERE date <= $1`
	_, err = env.DB.Exec(query, now)
	if err != nil {
		panic(err)
	}
}
func (env *Env) NearestFixture() Fixture {
	var fixture Fixture
	timeNow := time.Now()
	err := env.DB.QueryRow(
		"SELECT id, date, home_team, away_team FROM public.fixtures WHERE date > $1 ORDER BY date ASC", timeNow).Scan(
		&fixture.FixtureID, &fixture.EventDate, &fixture.HomeTeam.TeamName, &fixture.AwayTeam.TeamName)
	if err != nil {
		panic(err)
	}
	fixture.TimeTo = fixture.EventDate.Sub(timeNow)
	if fixture.HomeTeam.TeamName == "Chelsea" {
		fixture.HomeTeam.TeamName = "@Chelsea"
	}
	if fixture.AwayTeam.TeamName == "Chelsea" {
		fixture.AwayTeam.TeamName = "@Chelsea"
	}
	return fixture
}

func (env *Env) GetLineup(fixture Fixture) string {
	uri := "https://server1.api-football.com/lineups/" + strconv.Itoa(fixture.FixtureID)

	client := http.Client{}

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}
	headers := http.Header{
		keyHeader: []string{tokenAPIFoot},
	}
	request.Header = headers
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	lineup := Lineup{}
	err = json.Unmarshal(body, &lineup)
	if err != nil || len(lineup.API.LineUps.Chelsea.StartXI) == 0 || lineup.API.Results == 0 {
		log.Println("lineup result: ", lineup.API.Results)
		log.Println("lineup formation: ", lineup.API.LineUps.Chelsea.Formation)
		log.Println("will check atfer 50 sec")
		time.Sleep(time.Second * 50)
		env.GetLineup(fixture)
	}
	text := " *Match of the day:*\n"
	text += "*" + fixture.HomeTeam.TeamName + " - " + fixture.AwayTeam.TeamName + "*\n"
	text += "\n *Line-up (" + lineup.API.LineUps.Chelsea.Formation + "):*\n"
	for _, player := range lineup.API.LineUps.Chelsea.StartXI {
		text += player.Player + " (" + strconv.Itoa(player.Number) + ")" + "\n"
	}

	logrus.Info(text)
	text += "\n *Substitutes*:\n"
	for _, player := range lineup.API.LineUps.Chelsea.Substitutes {
		text += player.Player + " (" + strconv.Itoa(player.Number) + ")" + "\n"
	}
	text += "\n *Match starts at*:\n"
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

func (env *Env) SendPost(text string) {
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=Markdown",
		botToken, testChat, url.QueryEscape(text))
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	logrus.Info(string(body))
}

func (env *Env) SetUp(postedFixture *int) int {
	env.GetFixtures(*postedFixture)
	fixture := env.NearestFixture()
	logrus.Info("vs team: ", fixture.HomeTeam.TeamName, fixture.AwayTeam.TeamName, " id: ", fixture.FixtureID)
	logrus.Info("will send post after ", fixture.TimeTo-(time.Minute*55))
	time.Sleep(fixture.TimeTo - (time.Minute * 55))
	text := env.GetLineup(fixture)
	env.SendPost(text)
	// env.DeleteFixture(fixture)
	postedFixture = &fixture.FixtureID
	return *postedFixture
}

func (env *Env) StatusCheck() Status {
	uri := "https://server1.api-football.com/status"

	client := http.Client{}

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}
	headers := http.Header{
		keyHeader: []string{tokenAPIFoot},
	}
	request.Header = headers
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	status := Status{}
	err = json.Unmarshal(body, &status)

	if err != nil {
		log.Fatalln(err)
	}
	status.API.Status.RequestsLeft = status.API.Status.RequestsLimitDay - status.API.Status.Requests
	return status
}
func (env *Env) Status(w http.ResponseWriter, r *http.Request) {
	status := env.StatusCheck()
	reqLeft := strconv.Itoa(status.API.Status.RequestsLeft)
	_, err := io.WriteString(w, `{"Requests remaining":`+reqLeft+`}`)
	if err != nil {
		log.Fatalln(err)
	}
}
func (env *Env) DeleteFixture(fixture Fixture) {
	query := `
	DELETE FROM fixtures
	WHERE id = $1`
	_, err := env.DB.Exec(query, fixture.FixtureID)
	if err != nil {
		panic(err)
	}
	logrus.Info("fixture deleted: ", fixture.FixtureID)
}
func (env *Env) CreateTable(w http.ResponseWriter, r *http.Request) {
	query := `
	CREATE TABLE fixtures(
		id INTEGER PRIMARY KEY,
		league_id INTEGER,
		date TIMESTAMP WITH TIME ZONE,
		home_team VARCHAR (50),
		away_team VARCHAR (50)
	 );`
	_, err := env.DB.Exec(query)
	if err != nil {
		panic(err)
	}
}
func DetermineListenAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8000"
	}
	return ":" + port
}
func Updater() {
	for {
		logrus.Info("server is up")
		uri := "https://chelsea-score-bot.herokuapp.com/"
		client := http.Client{}
		request, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			log.Fatalln(err)
		}
		client.Do(request)
		time.Sleep(25 * time.Minute)
	}
}

func (env *Env) Start() {
	var (
		postedFixture int
	)
	for {
		postedFixture = env.SetUp(&postedFixture)
	}
}
