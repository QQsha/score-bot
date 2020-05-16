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
	botToken string
	chatID   string
}

func NewBotRepository(botToken, chatID string) *BotRepository {
	return &BotRepository{
		botToken: botToken,
		chatID:   chatID,
	}
}
func (r BotRepository) GetChatID() string {
	return r.chatID
}
func (r BotRepository) SendPost(text string) {
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=Markdown",
		r.botToken, r.chatID, url.QueryEscape(text))
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
	banUntil := time.Now().Add(time.Second * time.Duration(banDuration)).Unix()
	uri := fmt.Sprintf(
		"https://api.telegram.org/bot%s/restrictChatMember?chat_id=%s&user_id=%s&until_date=%s&permissions:{can_send_messages:false}",
		r.botToken, strconv.Itoa(-1001276457176), strconv.Itoa(userID), strconv.Itoa(int(banUntil)))
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
