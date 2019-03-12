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
	Me     *User
	token  string
	client *http.Client
}

// NewBot create a fresh Bot
func NewBot() *Bot {
	bot := &Bot{
		token:  BotToken,
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
	url := fmt.Sprintf(TelegramEndpoint, bot.token, method)
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

// ReceiveUpdates setups a Webhook and listens to it
func (bot *Bot) ReceiveUpdates(host string) (UpdateChannel, error) {
	pattern := fmt.Sprintf(WebhookURL, bot.token)
	_, err := bot.setupWebhook(host + pattern)
	if err != nil {
		return nil, err
	}
	return bot.listenToWebhook(pattern), nil
}

func (bot *Bot) setupWebhook(urlForUpdates string) (*TelegramResponse, error) {
	updatesToReceive := []string{"message", "edited_message", "channel_post", "edited_channel_post"}
	v := url.Values{}
	v.Add("url", urlForUpdates)
	v.Add("max_connections", strconv.Itoa(WebhookMaxConnections))
	v.Add("allowed_updates", strings.Join(updatesToReceive, ","))
	return bot.Request("setWebhook", v)
}

func (bot *Bot) listenToWebhook(pattern string) UpdateChannel {
	channel := make(chan Update, 100)

	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		var update Update
		json.Unmarshal(bytes, &update)

		channel <- update
	})

	return channel
}

// RemoveUpdates removes a Webhook
func (bot *Bot) RemoveUpdates() (*TelegramResponse, error) {
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
