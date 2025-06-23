package main

import (
	"log"
	"os"
	"tg-remote/internal/bot"
	"time"

	"github.com/dotenv-org/godotenvvault"
)

func main() {
	err := godotenvvault.Load()
	if err != nil {
		log.Print("no env")
	}

	// Get bot token from environment variable
	botToken := os.Getenv("TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// Initialize Telegram bot
	telegramBot := bot.NewTelegramBot(botToken)

	log.Println("🤖 Telegram bot started successfully!")
	log.Println("📡 Listening for updates...")

	// Main polling loop
	for {
		updates, err := telegramBot.GetUpdates()
		if err != nil {
			log.Printf("Error getting updates: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			if update.Message.Text != "" {
				log.Printf("Received message from %s: %s",
					update.Message.From.Username, update.Message.Text)

				// Handle message in goroutine (no jobicy parameter needed)
				go bot.HandleMsg(telegramBot, update.Message)
			}

			telegramBot.UpdateOffset(update.UpdateID)
		}

		// Small delay to prevent excessive API calls
		time.Sleep(1 * time.Second)
	}
}
