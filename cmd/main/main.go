package main

import (
	botapi "JustNotifierBot/internal/botapi"
	log "JustNotifierBot/internal/logger"
	"net/http"
)

func main() {
	bot := botapi.NewBot()

	updates, err := bot.ReceiveUpdates("127.0.0.1:80")
	if err != nil {
		log.Printf("ReceiveUpdates failed: %s", err.Error())
	}

	go http.ListenAndServe(":8080", nil)

	for update := range updates {
		log.Printf("Next update: %#v", update)
	}
}
