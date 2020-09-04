package posts

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/QQsha/score-bot/models"
)

const (
	timeFormat = "15:04"
)

func CreatePostLineup(fixture models.Fixture, lineup models.Lineup) string {
	text := "📣💙*Match of the day:*💙📣\n"
	text += "*" + fixture.HomeTeam.TeamName + " - " + fixture.AwayTeam.TeamName + "*\n"
	text += "\n *Line-up (" + lineup.API.LineUps.Chelsea.Formation + "):*\n"
	for _, player := range lineup.API.LineUps.Chelsea.StartXI {
		text += player.Player + " (" + strconv.Itoa(player.Number) + ")" + "\n"
	}
	text += "\n *Substitutes*:\n"
	for _, player := range lineup.API.LineUps.Chelsea.Substitutes {
		text += player.Player + " (" + strconv.Itoa(player.Number) + ")" + "\n"
	}
	text += "\n *Match starts today at*:\n"
	loc, _ := time.LoadLocation("Asia/Tehran")
	text += "🇮🇷*Tehran*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Africa/Lagos")
	text += "🇳🇬*Abuja*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Europe/Moscow")
	text += "🇷🇺*Moscow*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Asia/Tashkent")
	text += "🇺🇿*Tashkent*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Europe/London")
	text += "🇬🇧*London*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	return text
}

func CreatePostNextGame(fixture models.Fixture) string {

	text := "*" + fixture.LeagueName + "*\n"
	text += "*" + fixture.HomeTeam.TeamName + " - " + fixture.AwayTeam.TeamName + "* \n"
	text += "\n *Match will be on " + fixture.EventDate.In(time.UTC).Format("2 January") + " at*:\n"
	loc, _ := time.LoadLocation("Asia/Tehran")
	text += "🇮🇷*Tehran*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Africa/Addis_Ababa")
	text += "🇪🇹*Addis Ababa*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Africa/Lagos")
	text += "🇳🇬*Abuja*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Europe/Moscow")
	text += "🇷🇺*Moscow*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Asia/Tashkent")
	text += "🇺🇿*Tashkent*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Asia/Kolkata")
	text += "🇮🇳*New Delhi*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"
	loc, _ = time.LoadLocation("Europe/London")
	text += "🇬🇧*London*: " + fixture.EventDate.In(loc).Format(timeFormat) + "\n"

	timeTo := time.Until(fixture.EventDate)

	text += "\n *Time left: " + fmtDuration(timeTo) + "*"

	text += "\n\n You can predict score to this match with command / predict 3-1"
	return text
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d Hours, %02d Minutes", h, m)
}

func CreatePostEvent(fixture models.FixtureDetails, event models.Event) string {
	post := timeConv(event.Elapsed)
	switch event.Type {
	case "Goal":
		post += "⚽️Goal by " + event.Player + "(" + event.Assist + ")\n" + GetCurrentScore(fixture)
	case "Card":
		switch event.Detail {
		case "Yellow Card":
			post += "📒" + event.Player + " got a yellow card."
		case "Red Card":
			post += "📕" + event.Player + " got a red card."
		}
	case "subst":
		post += "🔄" + event.Player + " substitute " + event.Assist
	}
	return post
}

func timeConv(minute int) string {
	return strconv.Itoa(minute) + "' "
}

func CreateMVPPost(fixture models.FixtureDetails) (string, []string) {
	players := getBestPlayers(fixture)
	question := "🥇VOTE your MVP of the match:\n\n"
	question += GetCurrentScore(fixture) + "\n"
	return question, players
}
func getBestPlayers(fixture models.FixtureDetails) []string {
	players := fixture.API.Fixtures[0].Players
	playersResult := make([]string, 0)

	for i, player := range players {
		if player.TeamID == 49 {
			rankFloat, _ := strconv.ParseFloat(player.Rating, 64)
			players[i].FloatRaiting = rankFloat
		}
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].FloatRaiting > players[j].FloatRaiting
	})
	if len(players) >= 10 {
		for _, player := range players[:10] {
			playersResult = append(playersResult, player.PlayerName)
			fmt.Println(player.FloatRaiting)
		}
	}
	return playersResult
}

func CreateFullTimePost(fixture models.FixtureDetails) string {
	post := "📣*Match finished!\n\n" + fixture.API.Fixtures[0].League.Name + "*\n"
	post += GetCurrentScore(fixture) + "\n\n"
	post += GetGoals(fixture) + "\n"
	post += "KEEP THE BLUE FLAG FLYING HIGH!💙"
	return post

}

func CreateWinnersPost(winners []models.User) string {
	if len(winners) == 0 {
		return ""
	}
	post := "\n\nConcragutlaitions to chat members, who predicted correct score: \n"
	for _, winner := range winners {
		//[Andrei Cucuschin](tg://user?id=46731206)
		post += "[" + winner.FirstName + " " + winner.LastName + "](tg://user?id=" + strconv.Itoa(winner.ID) + ") \n"
	}

	return post
}

func CreateLeaderboardPost(winners []models.User) string {
	if len(winners) == 0 {
		return "*Prediction Leaderboard:* \n\n no win prediction yet"
	}
	post := "Prediction Leaderboard: \n"
	for i, winner := range winners {
		post += "\n*№" + strconv.Itoa(i+1) + ".* " + winner.FirstName + " " + winner.LastName + " (wins: " + strconv.Itoa(winner.Wins) + ")"
		for _, fixture := range winner.Fixtures {
			post += fixture + "\n"
		}
	}
	return post
}

func CreateStatsPost(fixture models.FixtureDetails) string {
	stats := fixture.API.Fixtures[0].Statistics
	post := "📊*Match statistics: *\n\n"
	post += GetCurrentScore(fixture) + "\n"
	post += centerlizer(fixture, true) + "*Ball posetion*\n"
	post += centerlizer(fixture, false) + stats.BallPossession.Home + " - " + stats.BallPossession.Away + "\n"
	post += centerlizer(fixture, true) + "*Total Shots*\n"
	post += centerlizer(fixture, false) + stats.TotalShots.Home + " - " + stats.TotalShots.Away + "\n"
	post += centerlizer(fixture, true) + "*Shots on target*\n"
	post += centerlizer(fixture, false) + stats.ShotsOnGoal.Home + " - " + stats.ShotsOnGoal.Away + "\n"
	post += centerlizer(fixture, true) + "*Shots off target*\n"
	post += centerlizer(fixture, false) + stats.ShotsOffGoal.Home + " - " + stats.ShotsOffGoal.Away + "\n"
	post += centerlizer(fixture, true) + "*Blocked shots*\n"
	post += centerlizer(fixture, false) + stats.BlockedShots.Home + " - " + stats.BlockedShots.Away + "\n"
	post += centerlizer(fixture, true) + "*Fouls*\n"
	post += centerlizer(fixture, false) + stats.Fouls.Home + " - " + stats.Fouls.Away + "\n"
	post += centerlizer(fixture, true) + "*Corner kicks*\n"
	post += centerlizer(fixture, false) + stats.CornerKicks.Home + " - " + stats.CornerKicks.Away + "\n"
	return post
}
func centerlizer(fixture models.FixtureDetails, title bool) string {
	var res string
	if !title {
		res += "\t\t\t\t"
	}
	for i := 0; i < len(fixture.API.Fixtures[0].HomeTeam.TeamName); i++ {
		res += "\t"
	}
	return res
}

func GetGoals(fixture models.FixtureDetails) string {
	var (
		result, secondTeamName, chelseaGoals, secondTeamGoals string
		chelseaGotGoals, secondGotGoals                       bool
	)
	for _, event := range fixture.API.Fixtures[0].Events {
		if event.Type == "Goal" {
			if event.TeamID == 49 {
				chelseaGoals += timeConv(event.Elapsed) + "⚽️" + event.Player + "(" + event.Assist + ")\n"
				chelseaGotGoals = true
				continue
			}
			secondTeamGoals += timeConv(event.Elapsed) + "⚽️" + event.Player + "(" + event.Assist + ")\n"
			secondGotGoals = true
			if secondTeamName == "" {
				secondTeamName = GetSeondTeamName(fixture)
			}
		}
	}
	if chelseaGotGoals {
		result += "@Chelsea goals:\n" + chelseaGoals
	}
	if secondGotGoals {
		result += "\n" + secondTeamName + " goals:\n" + secondTeamGoals
	}
	return result
}

func GetCurrentScore(fixture models.FixtureDetails) string {
	text := fixture.API.Fixtures[0].HomeTeam.TeamName + "   " + strconv.Itoa(fixture.API.Fixtures[0].GoalsHomeTeam) + " - " + strconv.Itoa(fixture.API.Fixtures[0].GoalsAwayTeam) + "   " + fixture.API.Fixtures[0].AwayTeam.TeamName + "\n"
	return text
}

func GetSeondTeamName(fixture models.FixtureDetails) string {
	if fixture.API.Fixtures[0].HomeTeam.TeamID == 49 {
		return fixture.API.Fixtures[0].AwayTeam.TeamName
	}
	return fixture.API.Fixtures[0].HomeTeam.TeamName
}

func GetRulePost() string {
	post := `
	*GROUP RULES:*
	 - *English only.*
	-* No links, mentions, and forwards allowed*, for keeping this chat clean from advertise.(auto Punishment: 24 hours in mute)
	- *Banned words:* (auto Punishment: 24 hours in mute)
	 "bets",
		"odds", 
		"fix", 
		"porn",
		"Arsenal".
	- *Be respectful* to other chat members and @Chelsea players.
	- *Racism or discrimination* will result in a ban.
	`
	return post
}

func GetEnglishPhrases() []string {
	phrases := []string{
		"dude, you are so cool, but can you speak english here pls",
		"Bro, have a nice day! But speak english in this chat pls",
		"Bro, i love you, but i dont understand you, just speak english",
		"You are so fucking cool my friend, but you should speak english in this chat!",
		"E N G L I S H",
		"Bro, english please",
		"Your language is so beautiful, but i dont understand it, you should switch to english.",
		"This is English? I dont think so.",
		"Or is it the name of a football player, or you do not speak English.",
	}
	return phrases
}

func GetArsenalPhrases() []string {
	phrases := []string{
		"Arsenal, what a losers",
		"Arsenal is so bad, im right?",
		"Arsanal",
		"stupid Arsenal",
		"wrost team in the word",
		"ars8nal",
		"The Best comedy team in EPL",
		"beeaaahh",
		"Arsneal in EPL, like Drinkwater in Chelsea.",
		"(2 + 2)*2 = Arsenal",
		"If Drinkwater will come in Arsenal, they will get plus one star in Fifa",
		"Good ebening",
		"Jesus, please help Arsenal to take 3 points",
		"Arsenal, they dont even can get 4th place",
		"i bet my medical licence, that Chelsea U-14 can beat Arsenal",
		"Arsenal, just lets they will be allowed to play in football with hands",
		"Arsenal dont have money for transfers, because they paying all their money to EA Sports, to keep their team in Fifa20",
		"82% of anual Arsenal's budget generates hotdog guys from Emirates",
		"Arsenal Best player of the Year - Anthony Taylor",
	}
	return phrases
}
func FantasyPost() string {
	post := "JOIN our Chelsea Chat Fantasy League - \n https://fantasy.premierleague.com/leagues/auto-join/hhmmkb"
	return post
}
