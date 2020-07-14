package mock

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/QQsha/score-bot/models"
	"github.com/QQsha/score-bot/repository"
)

type APIRepositoryMock struct {
}

func NewAPIRepositoryMock() *APIRepositoryMock {
	return &APIRepositoryMock{}
}

func (api APIRepositoryMock) GetFixtures() models.Fixtures {
	jsonFile, err := os.Open("repository/mock/Chelsea_Fixtures_mock.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	fixtures := models.Fixtures{}
	err = json.Unmarshal(byteValue, &fixtures)
	if err != nil {
		log.Fatalln(err)
	}
	return fixtures
}

func (api APIRepositoryMock) GetLineup(fixtureID int) models.Lineup {
	lineup := models.Lineup{}
	if fixtureID == 0 {
		return lineup
	}
	jsonFile, err := os.Open("repository/mock/Chelsea_Lineup_mock.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &lineup)
	if err != nil {
		log.Fatalln(err)
	}
	return lineup
}
func (api APIRepositoryMock) StatusCheck() repository.Status {
	status := repository.Status{}
	return status
}

func (api APIRepositoryMock) GetFixtureDetails(fixtureID int) models.FixtureDetails {
	jsonFile, err := os.Open("repository/mock/Chelsea_Fixture_detail_mock.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	fixtureDetails := models.FixtureDetails{}
	err = json.Unmarshal(byteValue, &fixtureDetails)
	if err != nil {
		log.Fatalln(err)
	}

	return fixtureDetails
}

func (api APIRepositoryMock) GetLeagues() models.Leagues {
	league := models.Leagues{}
	league.API.Leagues = append(league.API.Leagues, models.League{ID: 345, Name: "APL"})

	return league
}
