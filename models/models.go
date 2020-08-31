package models

import (
	"strconv"
	"time"
)

// https://server2.api-football.com/fixtures/team/49
type Fixtures struct {
	API struct {
		Results  int       `json:"results"`
		Fixtures []Fixture `json:"fixtures"`
	} `json:"api"`
}
type Leagues struct {
	API struct {
		Results int      `json:"results"`
		Leagues []League `json:"leagues"`
	} `json:"api"`
}
type League struct {
	ID   int    `json:"league_id"`
	Name string `json:"name"`
}
type Fixture struct {
	FixtureID       int `json:"fixture_id"`
	LeagueID        int `json:"league_id"`
	LeagueName      string
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

type FixtureEvents struct {
	API struct {
		Results int     `json:"results"`
		Events  []Event `json:"events"`
	} `json:"api"`
}

type FixtureDetails struct {
	API struct {
		Results  int `json:"results"`
		Fixtures []struct {
			FixtureID int `json:"fixture_id"`
			LeagueID  int `json:"league_id"`
			League    struct {
				Name    string `json:"name"`
				Country string `json:"country"`
				Logo    string `json:"logo"`
				Flag    string `json:"flag"`
			} `json:"league"`
			EventDate       time.Time `json:"event_date"`
			EventTimestamp  int       `json:"event_timestamp"`
			FirstHalfStart  int       `json:"firstHalfStart"`
			SecondHalfStart int       `json:"secondHalfStart"`
			Round           string    `json:"round"`
			Status          string    `json:"status"`
			StatusShort     string    `json:"statusShort"`
			Elapsed         int       `json:"elapsed"`
			Venue           string    `json:"venue"`
			Referee         string    `json:"referee"`
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
			Events  []Event `json:"events"`
			Lineups struct {
				Leicester struct {
					Coach     string `json:"coach"`
					CoachID   int    `json:"coach_id"`
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
				} `json:"Leicester"`
				Chelsea struct {
					Coach     string `json:"coach"`
					CoachID   int    `json:"coach_id"`
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
				} `json:"Chelsea"`
			} `json:"lineups"`
			Statistics struct {
				ShotsOnGoal struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Shots on Goal"`
				ShotsOffGoal struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Shots off Goal"`
				TotalShots struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Total Shots"`
				BlockedShots struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Blocked Shots"`
				ShotsInsidebox struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Shots insidebox"`
				ShotsOutsidebox struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Shots outsidebox"`
				Fouls struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Fouls"`
				CornerKicks struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Corner Kicks"`
				Offsides struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Offsides"`
				BallPossession struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Ball Possession"`
				YellowCards struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Yellow Cards"`
				RedCards struct {
					Home interface{} `json:"home"`
					Away interface{} `json:"away"`
				} `json:"Red Cards"`
				GoalkeeperSaves struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Goalkeeper Saves"`
				TotalPasses struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Total passes"`
				PassesAccurate struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Passes accurate"`
				Passes struct {
					Home string `json:"home"`
					Away string `json:"away"`
				} `json:"Passes %"`
			} `json:"statistics"`
			Players []struct {
				EventID       int    `json:"event_id"`
				UpdateAt      int    `json:"updateAt"`
				PlayerID      int    `json:"player_id"`
				PlayerName    string `json:"player_name"`
				TeamID        int    `json:"team_id"`
				TeamName      string `json:"team_name"`
				Number        int    `json:"number"`
				Position      string `json:"position"`
				Rating        string `json:"rating"`
				FloatRaiting  float64
				MinutesPlayed int         `json:"minutes_played"`
				Captain       string      `json:"captain"`
				Substitute    string      `json:"substitute"`
				Offsides      interface{} `json:"offsides"`
				Shots         struct {
					Total int `json:"total"`
					On    int `json:"on"`
				} `json:"shots"`
				Goals struct {
					Total    int `json:"total"`
					Conceded int `json:"conceded"`
					Assists  int `json:"assists"`
				} `json:"goals"`
				Passes struct {
					Total    int `json:"total"`
					Key      int `json:"key"`
					Accuracy int `json:"accuracy"`
				} `json:"passes"`
				Tackles struct {
					Total         int `json:"total"`
					Blocks        int `json:"blocks"`
					Interceptions int `json:"interceptions"`
				} `json:"tackles"`
				Duels struct {
					Total int `json:"total"`
					Won   int `json:"won"`
				} `json:"duels"`
				Dribbles struct {
					Attempts int `json:"attempts"`
					Success  int `json:"success"`
					Past     int `json:"past"`
				} `json:"dribbles"`
				Fouls struct {
					Drawn     int `json:"drawn"`
					Committed int `json:"committed"`
				} `json:"fouls"`
				Cards struct {
					Yellow int `json:"yellow"`
					Red    int `json:"red"`
				} `json:"cards"`
				Penalty struct {
					Won      int `json:"won"`
					Commited int `json:"commited"`
					Success  int `json:"success"`
					Missed   int `json:"missed"`
					Saved    int `json:"saved"`
				} `json:"penalty"`
			} `json:"players"`
		} `json:"fixtures"`
	} `json:"api"`
}

func (fixture FixtureDetails) FixtureScore() string {
	if len(fixture.API.Fixtures) == 0 {
		return ""
	}
	return fixture.API.Fixtures[0].HomeTeam.TeamName + " " + strconv.Itoa(fixture.API.Fixtures[0].GoalsHomeTeam) + " - " + strconv.Itoa(fixture.API.Fixtures[0].GoalsAwayTeam) + " " + fixture.API.Fixtures[0].AwayTeam.TeamName
}

type Event struct {
	Elapsed     int    `json:"elapsed"`
	ElapsedPlus int    `json:"elapsed_plus"`
	TeamID      int    `json:"team_id"`
	TeamName    string `json:"teamName"`
	PlayerID    int    `json:"player_id"`
	Player      string `json:"player"`
	AssistID    int    `json:"assist_id"`
	Assist      string `json:"assist"`
	Type        string `json:"type"`
	Detail      string `json:"detail"`
	Comments    string `json:"comments"`
}

type Updates struct {
	Ok       bool          `json:"ok"`
	Messages []FullMessage `json:"result"`
}

type FullMessage struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID                          int    `json:"id"`
			Title                       string `json:"title"`
			Type                        string `json:"type"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
		} `json:"chat"`
		ForwardFrom struct {
			ID       int    `json:"id"`
			Type     string `json:"type"`
			Username string `json:"username"`
		} `json:"forward_from,omitempty"`
		ForwardFromChat struct {
			ID       int    `json:"id"`
			Type     string `json:"type"`
			Username string `json:"username"`
		} `json:"forward_from_chat,omitempty"`
		ForwardDate     int        `json:"forward_date,omitempty"`
		Caption         string     `json:"caption,omitempty"`
		CaptionEntities []Entities `json:"caption_entities,omitempty"`
		Date            int        `json:"date"`
		Text            string     `json:"text"`
		Dice            struct {
			Emoji string `json:"emoji"`
		}
	} `json:"message"`
}

type Entities struct {
	Type string `json:"type"`
}
type Spam struct {
	Word        string `json:"word"`
	BanDuration int    `json:"ban"`
}

type ChatUser struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}
type User struct {
	ID        int      `json:"id"`
	IsBot     bool     `json:"is_bot"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Username  string   `json:"username"`
	Wins      int      `json:"-"`
	Fixtures  []string `json:"-"`
}
type Result struct {
	User   User   `json:"user"`
	Status string `json:"status"`
}
