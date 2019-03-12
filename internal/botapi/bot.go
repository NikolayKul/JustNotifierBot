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
)

// Bot represents the API
type Bot struct {
	Me     *User
	Token  string
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

// SetupWebhook setups a Webhook where Telegram sends Updates to
func (bot *Bot) SetupWebhook(urlForUpdates string, maxConnections int) (*TelegramResponse, error) {
	// updatesToReceive := []string{"message", "edited_message", "channel_post", "edited_channel_post"}
	v := url.Values{}
	v.Add("url", urlForUpdates)
	v.Add("max_connections", strconv.Itoa(maxConnections))
	// v.Add("allowed_updates", strings.Join(updatesToReceive, ","))
	return bot.request("setWebhook", v)
}

// RemoveWebhook removes a Webhook
func (bot *Bot) RemoveWebhook() (*TelegramResponse, error) {
	return bot.request("setWebhook", url.Values{})
}

func (bot *Bot) request(method string, params url.Values) (*TelegramResponse, error) {
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

func (bot *Bot) getMe() (*User, error) {
	resp, err := bot.request("getMe", nil)
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
