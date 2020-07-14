package main_test

import (
	"os"
	"strconv"

	"github.com/QQsha/score-bot/repository"
	mock "github.com/QQsha/score-bot/repository/mock"
	"github.com/QQsha/score-bot/usecase"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ScoreBot", func() {
	var antiSpamBot *usecase.FixtureUseCase
	BeforeEach(func() {
		postgresConn := os.Getenv("POSTGRES_SCOREBOT")
		chatID := strconv.Itoa(-1001276457176)
		antiSpamBotToken := "1238653037:AAE5szQlripWfNcHPVnBF2UasrDlZ-KS_nk"
		antiSpamBotRepo := repository.NewBotRepository(antiSpamBotToken, chatID)
		detectLangApi := repository.NewLanguageAPI()
		db, _, _ := repository.Connect(postgresConn)
		fixtureRepo := repository.NewFixtureRepository(db)
		apiRepo := mock.NewAPIRepositoryMock()
		antiSpamBot = usecase.NewFixtureUseCase(*fixtureRepo, *antiSpamBotRepo, *apiRepo, *detectLangApi)
	})
	Describe("GetLineup", func() {
		Context("when linup is ready", func() {
			It("should return Chelsea linup", func() {
				lineup := antiSpamBot.GetLineup(2333)
				Expect(len(lineup.API.LineUps.Chelsea.StartXI)).Should(Equal(11))
			})
		})
	})

})
