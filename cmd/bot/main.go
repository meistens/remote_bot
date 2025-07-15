package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"tg-remote/internal/bot"

	"github.com/dotenv-org/godotenvvault"
)

func main() {
	// Load environment variables
	err := godotenvvault.Load()
	if err != nil {
		log.Print("No .env file found, using system environment variables")
	}

	// Command line flags
	var (
		webhookURL = flag.String("url", "", "Webhook URL (e.g., https://your-app.vercel.app/api/webhook) or custom if you can afford it")
		delete     = flag.Bool("delete", false, "Delete existing webhook")
	)
	flag.Parse()

	// Get bot token from environment
	botToken := os.Getenv("TOKEN")
	if botToken == "" {
		log.Fatal("TOKEN environment variable is required")
	}

	// Create bot instance
	telegramBot := bot.NewTelegramBot(botToken)

	if *delete {
		// Delete webhook
		fmt.Println("Deleting webhook...")
		err := telegramBot.DeleteWebhook()
		if err != nil {
			log.Fatalf("Failed to delete webhook: %v", err)
		}
		fmt.Println("✅ Webhook deleted successfully!")
		return
	}

	if *webhookURL == "" {
		fmt.Println("Usage:")
		fmt.Println("  Set webhook:    go run cmd/setup/main.go -url https://your-app.vercel.app/api/webhook")
		fmt.Println("  Delete webhook: go run cmd/setup/main.go -delete")
		return
	}

	// Set webhook
	fmt.Printf("Setting webhook to: %s\n", *webhookURL)
	err = telegramBot.SetWebhook(*webhookURL)
	if err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}

	fmt.Println("✅ Webhook set successfully!")
	fmt.Println("Your bot is now ready to receive updates via webhook.")
}
