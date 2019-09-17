package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	keyHeader = "X-RapidAPI-Key"
	token     = "9128771ca86462be53b41393a002341e"
	ddbURI    = "postgres://xsqidgwwvwvgkm:1e82bd5c5b23996ee1ed11dfaa89447adc5c524999c574b6c24b67c0c1a22604@ec2-75-101-153-56.compute-1.amazonaws.com:5432/ddgva2m0b3akm5"
	botToken  = "840859313:AAFfNUxxiaw6MIj9_5XSIeelJv7gns8qRqk"
	chatID    = "@Chelsea"
	testChat  = "-1001279121498"
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
		} `json:"lineUps"`
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

func (env *Env) GetFixtures() {
	status := env.StatusCheck()
	now := time.Now()
	if status.API.Status.RequestsLeft == 0 {
		time.Sleep(time.Hour)
		env.GetFixtures()
	}
	uri := "https://server1.api-football.com/fixtures/team/49"

	client := http.Client{}

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}
	headers := http.Header{
		"X-RapidAPI-Key": []string{"9128771ca86462be53b41393a002341e"},
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
			if fixture.StatusShort == "NS" {
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
		WHERE date < $1`
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
		"X-RapidAPI-Key": []string{"9128771ca86462be53b41393a002341e"},
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
	// fmt.Println(string(body))
	lineup := Lineup{}
	err = json.Unmarshal(body, &lineup)
	// fmt.Println(lineup)
	if err != nil {
		log.Fatalln(err)
	}
	if lineup.API.Results == 0 {
		time.Sleep(time.Second * 40)
		env.GetLineup(fixture)
	}
	text := "*Starting match*:%0A"
	text += fixture.HomeTeam.TeamName + " - " + fixture.AwayTeam.TeamName + "%0A"
	text += "%0A*Line-up*: %0A"
	for _, player := range lineup.API.LineUps.Chelsea.StartXI {
		text += player.Player + " (" + strconv.Itoa(player.Number) + ") %0A"
	}
	text += "%0A*Substitutes*: %0A"
	for _, player := range lineup.API.LineUps.Chelsea.Substitutes {
		text += player.Player + " (" + strconv.Itoa(player.Number) + ") %0A"
	}
	text += "%0A*Match starts at*: %0A"
	loc, _ := time.LoadLocation("Asia/Tehran")
	text += "🇮🇷*Tehran*: " + fixture.EventDate.In(loc).Format("2006-01-02 15:04") + "%0A"
	loc, _ = time.LoadLocation("Africa/Lagos")
	text += "🇳🇬*Abuja*: " + fixture.EventDate.In(loc).Format("2006-01-02 15:04") + "%0A"
	loc, _ = time.LoadLocation("Europe/Moscow")
	text += "🇷🇺*Moscow*: " + fixture.EventDate.In(loc).Format("2006-01-02 15:04") + "%0A"
	loc, _ = time.LoadLocation("Asia/Almaty")
	text += "🇰🇿*Nur-Sultan*: " + fixture.EventDate.In(loc).Format("2006-01-02 15:04") + "%0A"
	loc, _ = time.LoadLocation("Europe/London")
	text += "🇬🇧*London*: " + fixture.EventDate.In(loc).Format("2006-01-02 15:04") + "%0A"
	return text
}

func (env *Env) SendPost(text string) {
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=Markdown",
		botToken, chatID, text)
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(body))
}

func (env *Env) SetUp() {
	env.GetFixtures()
	fixture := env.NearestFixture()
	fmt.Println("vs team: ", fixture.HomeTeam.TeamName, fixture.AwayTeam.TeamName)
	fmt.Println("will send post after ", fixture.TimeTo-(time.Minute*50))
	time.Sleep(fixture.TimeTo - (time.Minute * 50))
	text := env.GetLineup(fixture)
	env.SendPost(text)
	env.DeleteFixture(fixture)
	env.SetUp()
}

func (env *Env) StatusCheck() Status {
	uri := "https://server1.api-football.com/status"

	client := http.Client{}

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}
	headers := http.Header{
		"X-RapidAPI-Key": []string{"9128771ca86462be53b41393a002341e"},
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
	// fmt.Println(string(body))
	status := Status{}
	err = json.Unmarshal(body, &status)
	// fmt.Println(lineup)
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
