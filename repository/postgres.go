package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/QQsha/score-bot/models"
	"github.com/lib/pq"
)

type FixtureRepositoryInterface interface {
	SaveFixtures(fixtures models.Fixtures) error
	NearestFixture(notPosted bool) (models.Fixture, error)
	AddLeader(member models.User, fixtureResult string) error
	GetLeaderboard() ([]models.User, error)
	AddPrediction(fixtureID, user_id, home_predict, away_predict int) error
	GetWinnersID(fixtureDetail models.FixtureDetails) ([]int, error)
	AddSpamWord(word string, banDuration int) error
	AddLeague(leagues models.Leagues) error
	DeleteSpamWord(word string) error
	GetSpamWords() ([]models.Spam, error)
	FixturePosted(fixtureID int) error
	MessageUpdates(chatID string, updateID int) error
	DateFix(fixureID int)
	ZeroEventer(fixureID int)
	EventDecrementer(fixureID int)
	EventIncrementer(fixureID int)
	EventCounter(fixureID int) (int, error)
	GetLastUpdate(chatID string) (int, error)
	PredictionCounter(fixtureID int) (int, error)
}

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
func (r FixtureRepository) SaveFixtures(fixtures models.Fixtures) error {
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
				return err
			}
		}
	}
	return nil
}

func (r FixtureRepository) NearestFixture(notPosted bool) (models.Fixture, error) {
	var fixture models.Fixture
	timeNow := time.Now()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qr := psql.Select("fixtures.id, fixtures.date, fixtures.home_team, fixtures.away_team, leagues.name").
		From("fixtures").
		Join("leagues ON leagues.id = fixtures.league_id").
		OrderBy("fixtures.date ASC").
		Limit(1).
		Where(sq.Gt{"fixtures.date": timeNow})
	if notPosted {
		qr = qr.Where(sq.Eq{"fixtures.posted": false})
	}
	query, args, err := qr.ToSql()
	if err != nil {
		return fixture, err
	}
	err = r.db.QueryRow(query, args...).
		Scan(
			&fixture.FixtureID,
			&fixture.EventDate,
			&fixture.HomeTeam.TeamName,
			&fixture.AwayTeam.TeamName,
			&fixture.LeagueName)
	if err != nil {
		return fixture, err
	}
	fixture.TimeTo = fixture.EventDate.Sub(timeNow)
	if fixture.HomeTeam.TeamName == "Chelsea" {
		fixture.HomeTeam.TeamName = "@Chelsea"
	}
	if fixture.AwayTeam.TeamName == "Chelsea" {
		fixture.AwayTeam.TeamName = "@Chelsea"
	}
	return fixture, nil
}

func (r FixtureRepository) CreateTableFixtures() error {
	// query := `
	// DROP TABLE fixtures;`
	// _, err := r.db.Exec(query)
	// if err != nil {
	// 	panic(err)
	// }
	query := `
		CREATE TABLE fixtures(
			id INTEGER PRIMARY KEY,
			league_id INTEGER,
			date TIMESTAMP WITH TIME ZONE,
			home_team VARCHAR (50),
			away_team VARCHAR (50),
			posted BOOLEAN DEFAULT false,
			events Int DEFAULT 0
		 );`
	_, err := r.db.Exec(query)
	return err
}

func (r FixtureRepository) CreateTableLastUpdate() error {
	query := `
		CREATE TABLE updates(
			group_name VARCHAR (50) PRIMARY KEY,
			last_update INTEGER
		 );`
	_, err := r.db.Exec(query)
	return err
}
func (r FixtureRepository) CreateTableLeagues() error {
	query := `
		CREATE TABLE leagues(
			id INTEGER PRIMARY KEY,
			name VARCHAR (250)
		 );`
	_, err := r.db.Exec(query)
	return err
}

func (r FixtureRepository) CreateTableSpamBase() error {
	query := `
		CREATE TABLE spam_words(
			word VARCHAR (250) PRIMARY KEY,
			ban_duration INTEGER
		 );`
	_, err := r.db.Exec(query)
	return err
}

func (r FixtureRepository) CreateTableEvents() error {
	query := `
		CREATE TABLE events(
			word VARCHAR (250) PRIMARY KEY,
			ban_duration INTEGER
		 );`
	_, err := r.db.Exec(query)
	return err
}

func (r FixtureRepository) CreateTablePredictions() error {
	query := `
		CREATE TABLE predictions(
			id VARCHAR (250) PRIMARY KEY,
			fixture_id INTEGER,
			user_id INTEGER,
			home_predict INTEGER,
			away_predict INTEGER,
			update_date TIMESTAMP WITH TIME ZONE
		 );`
	_, err := r.db.Exec(query)
	return err
}

func (r FixtureRepository) CreateTableLeaderboard() error {
	query := `
	DROP TABLE leaderboard;`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	query = `
		CREATE TABLE leaderboard(
			user_id INTEGER PRIMARY KEY,
			wins INTEGER,
			username VARCHAR (250),
			first_name VARCHAR (250),
			last_name VARCHAR (250),
			fixtures TEXT[]
		 );`
	_, err = r.db.Exec(query)
	return err
}

func (r FixtureRepository) AddLeader(member models.User, fixtureResult string) error {
	arr := make([]string, 1)
	arr = append(arr, fixtureResult)
	query := `
	INSERT INTO leaderboard (user_id, username, first_name, last_name, wins, fixtures)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (user_id)
	DO
	UPDATE
	SET 
	username = EXCLUDED.username,
	first_name = EXCLUDED.first_name,
	last_name = EXCLUDED.last_name,
	wins = EXCLUDED.wins + 1,
	fixtures = array_append(EXCLUDED.fixtures, $7);`
	_, err := r.db.Exec(query, member.ID, member.Username, member.FirstName, member.LastName, 1, pq.Array(arr), fixtureResult)
	return err
}

func (r FixtureRepository) GetLeaderboard() ([]models.User, error) {
	query := `
	SELECT user_id, first_name, last_name, wins, fixtures
	FROM leaderboard
	ORDER BY wins DESC
	LIMIT 10`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	winners := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Wins, pq.Array(&user.Fixtures))
		if err != nil {
			return nil, err
		}
		winners = append(winners, user)
	}
	return winners, nil
}

func (r FixtureRepository) AddPrediction(fixtureID, user_id, home_predict, away_predict int) error {
	query := `
	INSERT INTO predictions (id, fixture_id, user_id, home_predict, away_predict, update_date)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id)
	DO
	UPDATE
	SET 
	update_date = EXCLUDED.update_date,
	home_predict = EXCLUDED.home_predict,
	away_predict = EXCLUDED.away_predict`
	predictID := strconv.Itoa(fixtureID) + "_" + strconv.Itoa(user_id)
	_, err := r.db.Exec(query, predictID, fixtureID, user_id, home_predict, away_predict, time.Now())
	return err
}

func (r FixtureRepository) GetWinnersID(fixtureDetail models.FixtureDetails) ([]int, error) {
	fmt.Println(fixtureDetail.API.Fixtures[0].FixtureID,
		fixtureDetail.API.Fixtures[0].GoalsHomeTeam,
		fixtureDetail.API.Fixtures[0].GoalsAwayTeam)
	query := `
	SELECT user_id
	FROM predictions
	WHERE fixture_id = $1
	AND home_predict = $2
	AND away_predict = $3`
	rows, err := r.db.Query(
		query,
		fixtureDetail.API.Fixtures[0].FixtureID,
		fixtureDetail.API.Fixtures[0].GoalsHomeTeam,
		fixtureDetail.API.Fixtures[0].GoalsAwayTeam)
	if err != nil {
		return nil, err
	}
	winners := make([]int, 0)
	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		winners = append(winners, userID)
	}
	return winners, nil
}

func (r FixtureRepository) AddSpamWord(word string, banDuration int) error {
	query := `
	INSERT INTO spam_words (word, ban_duration)
	VALUES ($1, $2)
	ON CONFLICT (word)
	DO NOTHING`
	_, err := r.db.Exec(query, word, banDuration)
	return err
}

func (r FixtureRepository) AddLeague(leagues models.Leagues) error {
	for _, league := range leagues.API.Leagues {
		query := `
		INSERT INTO leagues (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id)
		DO NOTHING`
		_, err := r.db.Exec(query, league.ID, league.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r FixtureRepository) DeleteSpamWord(word string) error {
	query := `
	DELETE FROM spam_words
  WHERE word = $1;`
	_, err := r.db.Exec(query, word)
	return err
}

func (r FixtureRepository) GetSpamWords() ([]models.Spam, error) {
	query := `
	SELECT word, ban_duration
	FROM spam_words`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	spamWords := make([]models.Spam, 0)
	for rows.Next() {
		var spam models.Spam
		err := rows.Scan(&spam.Word, &spam.BanDuration)
		if err != nil {
			return nil, err
		}
		spamWords = append(spamWords, spam)
	}
	return spamWords, nil
}

func (r FixtureRepository) FixturePosted(fixtureID int) error {
	query := `
	UPDATE fixtures
	SET posted = $1
	WHERE id = $2;`
	_, err := r.db.Exec(query, true, fixtureID)
	return err
}

func (r FixtureRepository) MessageUpdates(chatID string, updateID int) error {
	query := `
	INSERT INTO updates (group_name, last_update)
	VALUES ($1, $2)
	ON CONFLICT (group_name)
	DO
	UPDATE
	SET last_update = EXCLUDED.last_update`
	_, err := r.db.Exec(query, chatID, updateID+1)
	return err
}

func (r FixtureRepository) GetLastUpdate(chatID string) (int, error) {
	query := `
	SELECT  last_update
	 FROM updates 
	 WHERE group_name = $1;`
	var lastUpdateID int
	err := r.db.QueryRow(query, chatID).Scan(&lastUpdateID)
	if err != nil {
		return lastUpdateID, err
	}
	return lastUpdateID, nil
}

func (r FixtureRepository) EventCounter(fixureID int) (int, error) {
	var eventCount int
	err := r.db.QueryRow(
		"SELECT events FROM fixtures WHERE id = $1", fixureID).Scan(&eventCount)
	return eventCount, err
}

func (r FixtureRepository) EventIncrementer(fixureID int) {
	r.db.QueryRow("UPDATE fixtures SET events = events + 1 WHERE id = $1;", fixureID)
}

func (r FixtureRepository) EventDecrementer(fixureID int) {
	r.db.QueryRow("UPDATE fixtures SET events = events - 1 WHERE id = $1;", fixureID)
}

func (r FixtureRepository) ZeroEventer(fixureID int) {
	r.db.QueryRow("UPDATE fixtures SET events = 0 WHERE id = $1;", fixureID)
}

func (r FixtureRepository) DateFix(fixureID int) {
	date := time.Date(2020, time.Month(8), 1, 19, 30, 0, 0, time.UTC)
	r.db.QueryRow("UPDATE fixtures SET date = $1 WHERE id = $2;", date, fixureID)
}

func (r FixtureRepository) PredictionCounter(fixtureID int) (int, error) {
	query := `
	SELECT COUNT(*) AS number_of_pred
	FROM predictions where Fixture_id = $1`
	row := r.db.QueryRow(query, fixtureID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
