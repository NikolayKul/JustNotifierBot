package main

import (
	botapi "JustNotifierBot/internal/botapi"
	log "JustNotifierBot/internal/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	bot := botapi.NewBot()

	updatesRoute := fmt.Sprintf("/updates/%s", bot.Token)
	_, err := bot.SetupWebhook(WebhookHost+updatesRoute, WebhookMaxConnections)
	if err != nil {
		log.Fatal(err)
	}

	updates := listenToRoute(updatesRoute)

	go http.ListenAndServe(":80", nil)

	for update := range updates {
		log.Printf("Next update: %#v", update)
	}
}

func listenToRoute(route string) updateChannel {
	channel := make(chan botapi.Update, 100)

	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		var update botapi.Update
		json.Unmarshal(bytes, &update)

		channel <- update
	})

	return channel
}

type updateChannel <-chan botapi.Update
