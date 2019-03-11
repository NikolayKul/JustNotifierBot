package botapi

import (
	log "JustNotifierBot/internal/logger"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Bot represents the API
type Bot struct {
	Token  string
	Me     *User
	client *http.Client
}

// NewBot create a fresh Bot
func NewBot() *Bot {
	bot := &Bot{
		Token:  BotToken,
		client: &http.Client{},
	}

	user, err := bot.getMe()
	if err != nil {
		log.Printf("`getMe` failed: %s", err.Error())
	}

	bot.Me = user
	return bot
}

// Request an endpoint with some params
func (bot *Bot) Request(method string, params url.Values) (*TelegramResponse, error) {
	url := fmt.Sprintf(TelegramEndpoint, bot.Token, method)
	log.Printf("Request: %s", url)

	rawResp, err := bot.client.PostForm(url, params)
	if err != nil {
		log.Printf("Request failed: %s", err.Error())
		return nil, err
	}
	defer rawResp.Body.Close()

	tgResp, err := decodeResponse(rawResp.Body)
	if err != nil {
		log.Printf("Decode failed: %s", err.Error())
		return nil, err
	}
	return tgResp, tgResp.GetError()
}

// SetWebhook to the Telegram API
func (bot *Bot) SetWebhook(urlForUpdates string) (*TelegramResponse, error) {
	updatesToReceive := []string{"message", "edited_message", "channel_post", "edited_channel_post"}
	v := url.Values{}
	v.Add("url", urlForUpdates)
	v.Add("max_connections", strconv.Itoa(WebhookMaxConnections))
	v.Add("allowed_updates", strings.Join(updatesToReceive, ","))
	return bot.Request("setWebhook", v)
}

// RemoveWebhook if there was any
func (bot *Bot) RemoveWebhook() (*TelegramResponse, error) {
	return bot.Request("setWebhook", url.Values{})
}

func (bot *Bot) getMe() (*User, error) {
	resp, err := bot.Request("getMe", nil)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(resp.Result, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func decodeResponse(responseBody io.Reader) (*TelegramResponse, error) {
	data, err := ioutil.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}

	var response TelegramResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
