package mock

import (
	"github.com/QQsha/score-bot/models"
)

type BotRepositoryMock struct {
	SendPostFunc      func(text string, replyMessageID *int) error
	SendPollFunc      func(text string, options []string) error
	GetUpdatesFunc    func(updateID int) (models.Updates, error)
	RestrictUserFunc  func(userID, banDuration int) error
	DeleteMessageFunc func(messageID int) error
	GetChatUserFunc   func(userID int) (models.User, error)
	GetChatIDFunc     func() string
}

func NewBotRepositoryMock() *BotRepositoryMock {
	return &BotRepositoryMock{}
}

func (s BotRepositoryMock) SendPost(text string, replyMessageID *int) error {
	return nil
}
func (s BotRepositoryMock) SendPoll(text string, options []string) error {
	return nil
}

func (s BotRepositoryMock) GetUpdates(updateID int) (models.Updates, error) {
	return s.GetUpdatesFunc(updateID)
}

func (s BotRepositoryMock) RestrictUser(userID, banDuration int) error {
	return nil
}

func (s BotRepositoryMock) DeleteMessage(messageID int) error {
	return nil
}

func (s BotRepositoryMock) GetChatUser(userID int) (models.User, error) {
	return s.GetChatUserFunc(userID)
}

func (s BotRepositoryMock) GetChatID() string {
	return s.GetChatIDFunc()
}
