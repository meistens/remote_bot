package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TG Bot Conn. Pool
type TelegramBot struct {
	Token   string
	BaseURL string
	Offset  int
}

// Part of the TG conn. pool
func NewTelegramBot(token string) *TelegramBot {
	return &TelegramBot{
		Token:   token,
		BaseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
		Offset:  0,
	}
}

// Receive incoming updates using long polling
func (bot *TelegramBot) GetUpdates() ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=30", bot.BaseURL, bot.Offset)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var telegramResp TelegramResponse
	if err := json.Unmarshal(body, &telegramResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !telegramResp.OK {
		return nil, fmt.Errorf("telegram API returned error")
	}

	return telegramResp.Result, nil
}

// Send text message, on success, sent Message is returned
func SendMessage(bot *TelegramBot, chatID int64, text string) error {
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

// Offset counter
func (bot *TelegramBot) UpdateOffset(updateID int) {
	bot.Offset = updateID + 1
}
