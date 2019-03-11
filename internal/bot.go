package botapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Bot represents the API
type Bot struct {
	Token  string `json:"token"`
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
		log.Print(err)
	}
	log.Printf("%#v", user)

	return bot
}

// Request an endpoint with some params
func (bot *Bot) Request(method string, params url.Values) (*TelegramResponse, error) {
	url := fmt.Sprintf(TelegramEndpoint, bot.Token, method)

	log.Printf("Request -> %s", url)

	rawResp, err := bot.client.PostForm(url, params)
	if err != nil {
		return nil, err
	}
	defer rawResp.Body.Close()

	tgResp, err := decodeResponse(rawResp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Response -> %#v", tgResp)

	if !tgResp.Ok {
		return tgResp, errors.New("Unknown error")
	}
	return tgResp, nil
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

	log.Printf("GetMe -> %#v", user)
	return &user, nil
}
