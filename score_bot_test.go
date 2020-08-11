package main_test

import (
	"github.com/QQsha/score-bot/models"
	"github.com/QQsha/score-bot/repository"
	mock "github.com/QQsha/score-bot/repository/mock"
	"github.com/QQsha/score-bot/usecase"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ScoreBot", func() {
	var antiSpamBot *usecase.FixtureUseCase
	BeforeSuite(func() {
		antiSpamBotRepo := mock.NewBotRepositoryMock()
		antiSpamBotRepo.GetChatUserFunc = func(id int) (models.User, error) {
			user := models.User{}
			user.FirstName = "Andrei"
			user.LastName = "Cucuschin"
			user.ID = id
			return user, nil
		}
		detectLangApi := repository.NewLanguageAPI()
		fixtureRepo := mock.NewRepositoryMock()
		fixtureRepo.GetWinnersIDFunc = func(models.FixtureDetails) ([]int, error) {
			return []int{46731206}, nil
		}
		fixtureRepo.AddLeaderFunc = func(models.User, string) error {
			return nil
		}
		apiRepo := mock.NewAPIRepositoryMock()
		antiSpamBot = usecase.NewFixtureUseCase(*fixtureRepo, *antiSpamBotRepo, *apiRepo, *detectLangApi)
	})
	Describe("GetLineup", func() {
		Context("when lineup is ready", func() {
			It("should return 11 players", func() {
				lineup := antiSpamBot.GetLineup(2333, 1)
				Expect(len(lineup.API.LineUps.Chelsea.StartXI)).Should(Equal(11))
			})
		})
	})
	Describe("IsSpam", func() {
		Context("when message inlclude spam word", func() {
			msg := models.FullMessage{}
			msg.Message.Text = "arsenal is the best"
			spamWords := make([]models.Spam, 1)
			spamWords = append(spamWords, models.Spam{Word: "arsenal"})
			It("should return true", func() {
				isSpam, _ := antiSpamBot.IsSpam(msg, spamWords)
				Expect(isSpam).Should(Equal(true))
			})
		})
		Context("when message from admin with spam word", func() {
			msg := models.FullMessage{}
			msg.Message.From.Username = "qqshaa"
			msg.Message.Text = "arsenal is the best"
			spamWords := make([]models.Spam, 1)
			spamWords = append(spamWords, models.Spam{Word: "arsenal"})
			It("should return false", func() {
				isSpam, _ := antiSpamBot.IsSpam(msg, spamWords)
				Expect(isSpam).Should(Equal(false))
			})
		})
		Context("when message is forwarded from another channel", func() {
			msg := models.FullMessage{}
			msg.Message.ForwardFrom.ID = 2243434
			It("should return true", func() {
				isSpam, _ := antiSpamBot.IsSpam(msg, nil)
				Expect(isSpam).Should(Equal(true))
			})
		})
		Context("when message include link to another channel", func() {
			msg := models.FullMessage{}
			msg.Message.CaptionEntities = append(msg.Message.CaptionEntities, models.Entities{Type: "mention"})
			It("should return true", func() {
				isSpam, _ := antiSpamBot.IsSpam(msg, nil)
				Expect(isSpam).Should(Equal(true))
			})
		})
	})
	Describe("GetWinners", func() {
		Context("when someone predicted correct score of the match", func() {
			It("should return post with winners", func() {
				fixDetail := models.FixtureDetails{}
				winners := antiSpamBot.GetWinners(fixDetail)
				expect := "\n\nConcragutlaitions to chat members, who predicted correct score: \n[Andrei Cucuschin](tg://user?id=46731206) \n"
				Expect(winners).Should(Equal(expect))
			})
		})
	})
})
