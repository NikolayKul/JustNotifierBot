package botapi

import (
	"encoding/json"
	"fmt"
)

// TelegramResponse is a Telegram API response
type TelegramResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"errorCode"`
	Description string          `json:"description"`
}

// GetError if response contains an error
func (resp *TelegramResponse) GetError() error {
	if resp.Ok {
		return nil
	}
	return fmt.Errorf("Telegram error: (errorCode=%d, desc=%s)", resp.ErrorCode, resp.Description)
}

// User is a Telegram API user
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"` // optional
	UserName  string `json:"username"`  // optional
}

// Update is a response from GetUpdates
type Update struct {
	UpdateID          int      `json:"update_id"`
	Message           *Message `json:"message"`
	EditedMessage     *Message `json:"edited_message"`
	ChannelPost       *Message `json:"channel_post"`
	EditedChannelPost *Message `json:"edited_channel_post"`
}

// Message in a Chat
type Message struct {
	MessageID            int              `json:"message_id"`
	From                 *User            `json:"from"` // optional
	Date                 int              `json:"date"`
	Chat                 *Chat            `json:"chat"`
	ForwardFrom          *User            `json:"forward_from"`            // optional
	ForwardFromChat      *Chat            `json:"forward_from_chat"`       // optional
	ForwardFromMessageID int              `json:"forward_from_message_id"` // optional
	ForwardDate          int              `json:"forward_date"`            // optional
	ReplyToMessage       *Message         `json:"reply_to_message"`        // optional
	EditDate             int              `json:"edit_date"`               // optional
	Text                 string           `json:"text"`                    // optional
	Entities             *[]MessageEntity `json:"entities"`                // optional
}

// MessageEntity contains information about data in a Message
type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	URL    string `json:"url"`  // optional
	User   *User  `json:"user"` // optional
}

// Chat for a Message
type Chat struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`                 // optional
	UserName    string `json:"username"`              // optional
	FirstName   string `json:"first_name"`            // optional
	LastName    string `json:"last_name"`             // optional
	Description string `json:"description,omitempty"` // optional
}
