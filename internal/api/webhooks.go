package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tg-remote/internal/bot"
	"tg-remote/internal/bot/services"
	"tg-remote/internal/bot/utils"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// allow only POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// parse webhook update
	var update bot.Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Printf("Error parsing webhook update: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	botToken := os.Getenv("TOKEN")
	if botToken == "" {
		log.Print("Bot token not found in environment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create bot instance
	telegramBot := bot.NewTelegramBot(botToken)

	// if a message exists
	if update.Message.Text != "" {
		log.Printf("Received message from %s: %s",
			update.Message.From.Username, update.Message.Text)

		// sync. handling of message
		err := handleMessage(telegramBot, update.Message)
		if err != nil {
			log.Printf("Error handling message: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// return statok
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handleMessage handles incoming messages sync.
func handleMessage(bot *bot.TelegramBot, message bot.Message) error {
	command, args := utils.ParseCommand(message.Text)

	switch command {
	case "/start":
		welcomeMsg := `Welcome to Remote Jobs Bot!

<b>Available Commands:</b>
• /jobs - Get latest remote jobs
• /jobs --count 5 - Get specific number of jobs
• /jobs --geo USA - Filter by location
• /jobs --industry tech - Filter by industry
• /jobs --tag python - Filter by technology tag
• /help - Show this help message

<b>Example:</b> <code>/jobs --count 3 --geo USA --tag golang</code>`

		return bot.SendMessage(message.Chat.ID, welcomeMsg)

	case "/help":
		helpMsg := `<b>Available Commands:</b>

<b>/jobs</b> - Get latest remote jobs (default: 5 jobs)

<b>Options:</b>
• <code>--count N</code> - Number of jobs (1-50)
• <code>--geo LOCATION</code> - Filter by location (USA, Europe, etc.)
• <code>--industry INDUSTRY</code> - Filter by industry
• <code>--tag TECHNOLOGY</code> - Filter by technology (python, golang, etc.)

<b>Examples:</b>
• <code>/jobs</code>
• <code>/jobs --count 10</code>
• <code>/jobs --geo USA --tag golang</code>
• <code>/jobs --industry tech --count 3</code>`

		return bot.SendMessage(message.Chat.ID, helpMsg)

	case "/jobs":
		return handleJobsCommand(bot, message, args)

	default:
		return bot.SendMessage(message.Chat.ID, "Unknown command. Type /help for available commands.")
	}
}

// handleJobsCommand handles the /jobs command
func handleJobsCommand(bot *bot.TelegramBot, message bot.Message, args map[string]string) error {
	// parse args
	count := 5
	if countStr, ok := args["count"]; ok {
		if c, err := strconv.Atoi(countStr); err == nil && c > 0 && c <= 50 {
			count = c
		}
	}

	geo := args["geo"]
	industry := args["industry"]
	tag := args["tag"]

	// send searching msg
	err := bot.SendMessage(message.Chat.ID, "🔍 Searching for jobs...")
	if err != nil {
		return err
	}

	// Jobicy client
	jobicyClient := services.NewJobicyClient()

	// Fetch jobs from API
	jobResp, err := services.GetJobs(jobicyClient, count, geo, industry, tag)
	if err != nil {
		log.Printf("Error fetching jobs: %v", err)
		return bot.SendMessage(message.Chat.ID, "❌ Sorry, I couldn't fetch jobs right now. Please try again later.")
	}

	if len(jobResp.Jobs) == 0 {
		return bot.SendMessage(message.Chat.ID, "❌ No jobs found with the specified criteria. Try different filters.")
	}

	// Send jobs (split into multiple messages if needed)
	var messageBuffer strings.Builder
	messageBuffer.WriteString(fmt.Sprintf("✅ Found %d remote jobs:\n\n", len(jobResp.Jobs)))

	for i, job := range jobResp.Jobs {
		jobMsg := services.FormatJobMsg(job)

		// Check if adding this job would exceed Telegram's message limit (4096 characters)
		if messageBuffer.Len()+len(jobMsg) > 4000 {
			// Send current buffer
			err := bot.SendMessage(message.Chat.ID, messageBuffer.String())
			if err != nil {
				return err
			}

			// Start new buffer
			messageBuffer.Reset()
			messageBuffer.WriteString(fmt.Sprintf("📋 Continuing jobs list (%d/%d):\n\n", i+1, len(jobResp.Jobs)))
		}

		messageBuffer.WriteString(jobMsg)
	}

	// Send remaining jobs
	if messageBuffer.Len() > 0 {
		return bot.SendMessage(message.Chat.ID, messageBuffer.String())
	}

	return nil
}
