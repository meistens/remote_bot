package services

import (
	"tg-remote/internal/bot"
)

// Jobicy
func NewJobicyClient() *bot.Jobicy {
	return &bot.Jobicy{
		BaseURL: "http://jobicy.com/api/v2/remote-jobs",
	}
}
