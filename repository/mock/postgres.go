package mock

import "github.com/QQsha/score-bot/models"

type RepositoryMock struct {
	NearestFixtureFunc   func(notPosted bool) (models.Fixture, error)
	GetWinnersIDFunc     func(fixtureDetail models.FixtureDetails) ([]int, error)
	SaveFixturesFunc     func(fixtures models.Fixtures) error
	AddLeaderFunc        func(member models.User, fixtureResult string) error
	GetLeaderboardFunc   func() ([]models.User, error)
	AddPredictionFunc    func(fixtureID, user_id, home_predict, away_predict int) error
	AddSpamWordFunc      func(word string, banDuration int) error
	AddLeagueFunc        func(leagues models.Leagues) error
	DeleteSpamWordFunc   func(word string) error
	GetSpamWordsFunc     func() ([]models.Spam, error)
	FixturePostedFunc    func(fixtureID int) error
	MessageUpdatesFunc   func(chatID string, updateID int) error
	DateFixFunc          func(fixureID int) error
	ZeroEventerFunc      func(fixureID int) error
	EventDecrementerFunc func(fixureID int) error
	EventIncrementerFunc func(fixureID int) error
	EventCounterFunc     func(fixureID int) (int, error)
	GetLastUpdateFunc    func(chatID string) (int, error)
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{}
}

func (s RepositoryMock) SaveFixtures(fixtures models.Fixtures) error {
	return s.SaveFixturesFunc(fixtures)
}

func (s RepositoryMock) NearestFixture(notPosted bool) (models.Fixture, error) {
	return s.NearestFixtureFunc(notPosted)
}

func (s RepositoryMock) GetWinnersID(fixtureDetail models.FixtureDetails) ([]int, error) {
	return s.GetWinnersIDFunc(fixtureDetail)
}

func (s RepositoryMock) AddLeader(member models.User, fixtureResult string) error {
	return s.AddLeaderFunc(member, fixtureResult)
}

func (s RepositoryMock) GetLeaderboard() ([]models.User, error) {
	return s.GetLeaderboardFunc()
}

func (s RepositoryMock) AddPrediction(fixtureID, user_id, home_predict, away_predict int) error {
	return s.AddPredictionFunc(fixtureID, user_id, home_predict, away_predict)
}

func (s RepositoryMock) AddSpamWord(word string, banDuration int) error {
	return s.AddSpamWordFunc(word, banDuration)
}
func (s RepositoryMock) AddLeague(leagues models.Leagues) error {
	return s.AddLeagueFunc(leagues)
}
func (s RepositoryMock) DeleteSpamWord(word string) error {
	return s.DeleteSpamWordFunc(word)
}
func (s RepositoryMock) GetSpamWords() ([]models.Spam, error) {
	return s.GetSpamWordsFunc()
}
func (s RepositoryMock) FixturePosted(fixtureID int) error {
	return s.FixturePostedFunc(fixtureID)
}
func (s RepositoryMock) MessageUpdates(chatID string, updateID int) error {
	return s.MessageUpdatesFunc(chatID, updateID)
}
func (s RepositoryMock) DateFix(fixureID int) {

}
func (s RepositoryMock) ZeroEventer(fixureID int) {

}
func (s RepositoryMock) EventDecrementer(fixureID int) {

}
func (s RepositoryMock) EventIncrementer(fixureID int) {

}
func (s RepositoryMock) EventCounter(fixureID int) (int, error) {
	return s.EventCounterFunc(fixureID)
}
func (s RepositoryMock) GetLastUpdate(chatID string) (int, error) {
	return s.GetLastUpdateFunc(chatID)
}
