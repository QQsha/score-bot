package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/QQsha/score-bot/models"
)

type FixtureRepository struct {
	db *sql.DB
}

func NewFixtureRepository(db *sql.DB) *FixtureRepository {
	return &FixtureRepository{
		db: db,
	}
}

func Connect(connection string) (*sql.DB, func() error, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, nil, err
	}
	return db, db.Close, nil
}
func (r FixtureRepository) SaveFixtures(fixtures models.Fixtures) {
	for _, fixture := range fixtures.API.Fixtures {
		if fixture.StatusShort == "NS" {
			query := `
					INSERT INTO fixtures (id, league_id, date, home_team, away_team)
					VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT (id)
					DO
					UPDATE
					SET date = EXCLUDED.date`
			_, err := r.db.Exec(
				query, fixture.FixtureID, fixture.LeagueID, fixture.EventDate, fixture.HomeTeam.TeamName, fixture.AwayTeam.TeamName)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (r FixtureRepository) NearestFixture() models.Fixture {
	var fixture models.Fixture
	timeNow := time.Now()
	err := r.db.QueryRow(
		"SELECT id, date, home_team, away_team FROM public.fixtures WHERE date > $1 and posted = $2 ORDER BY date ASC", timeNow, false).Scan(
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

func (r FixtureRepository) CreateTableFixtures() {
	query := `
	DROP TABLE fixtures;`
	_, err := r.db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `
		CREATE TABLE fixtures(
			id INTEGER PRIMARY KEY,
			league_id INTEGER,
			date TIMESTAMP WITH TIME ZONE,
			home_team VARCHAR (50),
			away_team VARCHAR (50),
			posted BOOLEAN DEFAULT false,
			events Int DEFAULT 0
		 );`
	_, err = r.db.Exec(query)
	if err != nil {
		panic(err)
	}

}

func (r FixtureRepository) CreateTableLastUpdate() {
	query := `
		CREATE TABLE updates(
			group_name VARCHAR (50) PRIMARY KEY,
			last_update INTEGER
		 );`
	_, err := r.db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (r FixtureRepository) CreateTableSpamBase() {
	query := `
		CREATE TABLE spam_words(
			word VARCHAR (250) PRIMARY KEY,
			ban_duration INTEGER
		 );`
	_, err := r.db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (r FixtureRepository) AddSpamWord(word string, banDuration int) {
	query := `
	INSERT INTO spam_words (word, ban_duration)
	VALUES ($1, $2)
	ON CONFLICT (word)
	DO NOTHING`
	_, err := r.db.Exec(query, word, banDuration)
	if err != nil {
		panic(err)
	}
}

func (r FixtureRepository) DeleteSpamWord(word string) {
	query := `
	DELETE FROM spam_words
  WHERE word = $1;`
	_, err := r.db.Exec(query, word)
	if err != nil {
		panic(err)
	}
}

func (r FixtureRepository) GetSpamWords() []models.Spam {
	query := `
	SELECT word, ban_duration
	FROM spam_words`
	rows, err := r.db.Query(query)
	if err != nil {
		panic(err)
	}
	spamWords := make([]models.Spam, 0)
	for rows.Next() {
		var spam models.Spam
		err := rows.Scan(&spam.Word, &spam.BanDuration)
		if err != nil {
			panic(err)
		}
		spamWords = append(spamWords, spam)
	}
	return spamWords
}

func (r FixtureRepository) FixturePosted(fixtureID int) {
	query := `
	UPDATE fixtures
	SET posted = $1
	WHERE id = $2;`
	_, err := r.db.Exec(query, true, fixtureID)
	if err != nil {
		panic(err)
	}
}

func (r FixtureRepository) MessageUpdates(chatID string, updateID int) {
	query := `
	INSERT INTO updates (group_name, last_update)
	VALUES ($1, $2)
	ON CONFLICT (group_name)
	DO
	UPDATE
	SET last_update = EXCLUDED.last_update`
	_, err := r.db.Exec(query, chatID, updateID+1)
	if err != nil {
		panic(err)
	}
}

func (r FixtureRepository) GetLastUpdate(chatID string) int {
	query := `
	SELECT  last_update
	 FROM updates 
	 WHERE group_name = $1;`
	var lastUpdateID int
	err := r.db.QueryRow(query, chatID).Scan(&lastUpdateID)
	if err != nil {
		fmt.Println(err)
	}
	return lastUpdateID
}

func (r FixtureRepository) EventCounter(fixureID int) int {
	var eventCount int
	err := r.db.QueryRow(
		"SELECT events FROM fixtures WHERE id = $1", fixureID).Scan(&eventCount)
	if err != nil {
		panic(err)
	}
	return eventCount
}
func (r FixtureRepository) EventIncrementer(fixureID int) {
	err := r.db.QueryRow(
		"UPDATE fixtures SET events = events + 1 WHERE id = $1;", fixureID)
	if err != nil {
		// fmt.Println(err)
	}
	return
}
