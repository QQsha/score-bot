package repository

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/QQsha/score-bot/models"
)

type APIRepository struct {
	token string
}
type APIRepositoryInterface interface {
	GetFixtures() models.Fixtures
	GetLineup(fixtureID int) models.Lineup
	StatusCheck() Status
	GetFixtureDetails(fixtureID int) models.FixtureDetails
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

func NewAPIRepository(token string) *APIRepository {
	return &APIRepository{
		token: token,
	}
}

const (
	keyHeader   = "X-RapidAPI-Key"
	fixturesURI = "http://v2.api-football.com/fixtures/team/49/next/100"
	lineupURI   = "http://v2.api-football.com/lineups/"
	statusURI   = "http://v2.api-football.com/status"
	eventsURI   = "http://v2.api-football.com/events/"
	fixtureURI  = "http://v2.api-football.com/fixtures/id/"
)

func (api APIRepository) GetFixtures() models.Fixtures {
	client := http.Client{}

	request, err := http.NewRequest("GET", fixturesURI, nil)
	if err != nil {
		log.Fatalln(err)
	}
	headers := http.Header{
		keyHeader: []string{api.token},
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
	fixtures := models.Fixtures{}
	err = json.Unmarshal(body, &fixtures)
	if err != nil {
		log.Fatalln(err)
	}
	return fixtures
}

func (api APIRepository) GetLineup(fixtureID int) models.Lineup {
	uri := lineupURI + strconv.Itoa(fixtureID)
	client := http.Client{}
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}
	headers := http.Header{
		keyHeader: []string{api.token},
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
	lineup := models.Lineup{}
	err = json.Unmarshal(body, &lineup)
	if err != nil {
		log.Fatalln(err)
	}
	return lineup
}
func (api APIRepository) StatusCheck() Status {
	client := http.Client{}
	request, err := http.NewRequest("GET", statusURI, nil)
	if err != nil {
		log.Fatalln(err)
	}
	headers := http.Header{
		keyHeader: []string{api.token},
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

func (api APIRepository) GetEvents(fixtureID, eventCount int) []models.Event {
	events := make([]models.Event, 0)
	uri := eventsURI + strconv.Itoa(fixtureID)
	client := http.Client{}
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}
	headers := http.Header{
		keyHeader: []string{api.token},
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
	fixtureEvents := models.FixtureEvents{}
	err = json.Unmarshal(body, &fixtureEvents)
	if err != nil {
		log.Fatalln(err)
	}
	if fixtureEvents.API.Results > eventCount {
		events = fixtureEvents.API.Events[eventCount:]
	}
	return events
}

func (api APIRepository) GetFixtureDetails(fixtureID int) models.FixtureDetails {
	uri := fixtureURI + strconv.Itoa(fixtureID)
	client := http.Client{}
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}
	headers := http.Header{
		keyHeader: []string{api.token},
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
	fixtureDetails := models.FixtureDetails{}
	err = json.Unmarshal(body, &fixtureDetails)
	if err != nil {
		log.Fatalln(err)
	}
	return fixtureDetails
}
