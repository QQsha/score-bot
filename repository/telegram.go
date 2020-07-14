package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/QQsha/score-bot/models"
)

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

func (r BotRepository) SendPost(text string, replyMessageID *int) {
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=Markdown&disable_notification=True",
		r.botToken, r.channelID, url.QueryEscape(text))
	if replyMessageID != nil {
		uri += "&reply_to_message_id=" + strconv.Itoa(*replyMessageID)
	}
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err, string(body))
	}
}

type MVP struct {
	Playeres []string `json:"options"`
}

func (r BotRepository) SendPoll(text string, options []string) {
	fmt.Println(len(options))
	mvp := MVP{}
	mvp.Playeres = options
	jsON, err := json.Marshal(mvp.Playeres)
	if err != nil {
		fmt.Println(err)
	}
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendPoll?chat_id=%s&question=%s&options=%s&disable_notification=True",
		r.botToken, r.channelID, url.QueryEscape(text),url.QueryEscape(string(jsON)) )
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err, string(body))
	}
}

func (r BotRepository) GetUpdates(updateID int) models.Updates {
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getUpdates?offset=%s",
		r.botToken, strconv.Itoa(updateID))
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err, string(body))
	}
	updates := models.Updates{}
	err = json.Unmarshal(body, &updates)
	if err != nil {
		log.Fatalln(err, string(body))
	}
	return updates
}

func (r BotRepository) RestrictUser(userID, banDuration int) {
	banUntil := time.Now().Add(time.Second * time.Duration(banDuration*3600)).Unix()
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/restrictChatMember?chat_id=%s&user_id=%s&until_date=%s&permissions:{can_send_messages:false}",
		r.botToken, r.channelID, strconv.Itoa(userID), strconv.Itoa(int(banUntil)))
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		log.Fatalln(err, string(body))
	}
}

func (r BotRepository) DeleteMessage(messageID int) {

	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/deleteMessage?chat_id=%s&message_id=%s",
		r.botToken, r.channelID, strconv.Itoa(messageID))
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		log.Fatalln(err, string(body))
	}
}
