package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TG Bot Connection Pool
type TelegramBot struct {
	Token   string
	BaseURL string
}

// Create new Telegram bot instance
func NewTelegramBot(token string) *TelegramBot {
	return &TelegramBot{
		Token:   token,
		BaseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
	}
}

// SendMessage sends a text message to the specified chat
func (bot *TelegramBot) SendMessage(chatID int64, text string) error {
	payload := SendMessagePayload{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/sendMessage", bot.BaseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %s", string(body))
	}
	return nil
}

// SetWebhook sets the webhook URL for the bot
func (bot *TelegramBot) SetWebhook(webhookURL string) error {
	payload := map[string]string{
		"url": webhookURL,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	url := fmt.Sprintf("%s/setWebhook", bot.BaseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook setup error: %s", string(body))
	}
	return nil
}

// DeleteWebhook removes the webhook
func (bot *TelegramBot) DeleteWebhook() error {
	url := fmt.Sprintf("%s/deleteWebhook", bot.BaseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook deletion error: %s", string(body))
	}
	return nil
}
