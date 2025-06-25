package bot

// Telegram Bot API structure

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"` // extends...
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from"` // extends...
	Chat      Chat   `json:"chat"` // extends...
	Date      int64  `json:"date"`
	Text      string `json:"text"`
}

type User struct {
	ID        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type TelegramResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type SendMessagePayload struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}
