package main

import (
	botapi "JustNotifierBot/internal/botapi"
	log "JustNotifierBot/internal/logger"
)

func main() {
	bot := botapi.NewBot()
	log.Printf("%#v", bot)
}
