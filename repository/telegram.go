package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/QQsha/score-bot/models"
)

type BotRepositoryInterface interface {
	SendPost(text string, replyMessageID *int) error
	SendPoll(text string, options []string) error 
	GetUpdates(updateID int) (models.Updates, error)
	RestrictUser(userID, banDuration int) error
	DeleteMessage(messageID int) error
	GetChatUser(userID int) (models.User, error)
	GetChatID() string
}
type BotRepository struct {
	botToken  string
	channelID string
}

func NewBotRepository(botToken, channelID string) *BotRepository {
	return &BotRepository{
		botToken:  botToken,
		channelID: channelID,
	}
}
func (r BotRepository) GetChatID() string {
	return r.channelID
}

func (r BotRepository) SendPost(text string, replyMessageID *int) error {
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=Markdown&disable_notification=True",
		r.botToken, r.channelID, url.QueryEscape(text))
	if replyMessageID != nil {
		uri += "&reply_to_message_id=" + strconv.Itoa(*replyMessageID)
	}
	_, err := http.Get(uri)

	return err

}

type MVP struct {
	Playeres []string `json:"options"`
}

func (r BotRepository) SendPoll(text string, options []string) error {
	fmt.Println(len(options))
	mvp := MVP{}
	mvp.Playeres = options
	jsON, err := json.Marshal(mvp.Playeres)
	if err != nil {
		return err
	}
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendPoll?chat_id=%s&question=%s&options=%s&disable_notification=True",
		r.botToken, r.channelID, url.QueryEscape(text), url.QueryEscape(string(jsON)))
	_, err = http.Get(uri)
	return err
}

func (r BotRepository) GetUpdates(updateID int) (models.Updates, error) {
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getUpdates?offset=%s",
		r.botToken, strconv.Itoa(updateID))
	resp, err := http.Get(uri)
	if err != nil {
		return models.Updates{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Updates{}, fmt.Errorf("cant read update %v, telegram response %v", err, string(body))
	}
	updates := models.Updates{}
	err = json.Unmarshal(body, &updates)
	if err != nil {
		return updates, err
	}
	return updates, nil
}

func (r BotRepository) RestrictUser(userID, banDuration int) error {
	banUntil := time.Now().Add(time.Second * time.Duration(banDuration*3600)).Unix()
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/restrictChatMember?chat_id=%s&user_id=%s&until_date=%s&permissions:{can_send_messages:false}",
		r.botToken, r.channelID, strconv.Itoa(userID), strconv.Itoa(int(banUntil)))
	_, err := http.Get(uri)
	if err != nil {
		return err
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (r BotRepository) DeleteMessage(messageID int) error {
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/deleteMessage?chat_id=%s&message_id=%s",
		r.botToken, r.channelID, strconv.Itoa(messageID))
	_, err := http.Get(uri)
	if err != nil {
		return err
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err, string(body))
	// }
	return nil
}

func (r BotRepository) GetChatUser(userID int) (models.User, error) {
	user := models.ChatUser{}
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getChatMember?chat_id=%s&user_id=%s",
		r.botToken, r.channelID, strconv.Itoa(userID))
	resp, err := http.Get(uri)
	if err != nil {
		return user.Result.User, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user.Result.User, err
	}
	fmt.Println(string(body))

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user.Result.User, err
	}
	return user.Result.User, nil
}
