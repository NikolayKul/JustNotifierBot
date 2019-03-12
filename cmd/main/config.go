package main

import "os"

const (
	// WebhookMaxConnections to set up webhook with
	WebhookMaxConnections = 40
)

// WebhookHost is a url for a webhook
var WebhookHost = os.Getenv("JUST_NOTIFIER_BOT_HOST")
