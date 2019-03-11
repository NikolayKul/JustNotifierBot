package botapi

import "os"

const (
	// TelegramEndpoint should be formatted to replace Token & Method
	TelegramEndpoint = "https://api.telegram.org/bot%s/%s"
)

// BotToken is a token retrieved from the BotFather
var BotToken = os.Getenv("JUST_NOTIFIER_BOT_TOKEN")
