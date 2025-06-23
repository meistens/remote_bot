package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"tg-remote/internal/bot/services"
	"tg-remote/internal/bot/utils"
)

// handleMsg handles incoming messages
func HandleMsg(bot *TelegramBot, message Message) {
	command, args := utils.ParseCommand(message.Text)

	switch command {
	case "/start":
		welcomeMsg := `welcome to remote jobs bot

		Commands:
		/jobs - Get latest remote jobs
		/jobs --count 5 - Get specific number of jobs
		/jobs --geo USA - Filter by location
		/jobs --industry tech - Filter by industry
		/jobs --tag python - Filter by technology tag
		/help - Show this help message

		Example: /jobs --count 3 --geo USA --tag golang`

		SendMessage(bot, message.Chat.ID, welcomeMsg)

	case "/help":
		helpMsg := `📚 Available Commands:

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

		SendMessage(bot, message.Chat.ID, helpMsg)

	case "/jobs":
		// Parse arguments
		count := 5 // default
		if countStr, ok := args["count"]; ok {
			if c, err := strconv.Atoi(countStr); err == nil && c > 0 && c <= 50 {
				count = c
			}
		}

		geo := args["geo"]
		industry := args["industry"]
		tag := args["tag"]

		// send searching msg
		SendMessage(bot, message.Chat.ID, "searching for jobs.......")

		// Create Jobicy client
		jobicyClient := services.NewJobicyClient()

		// fetch jobs from API
		jobResp, err := services.GetJobs(jobicyClient, count, geo, industry, tag)
		if err != nil {
			log.Printf("Error fetching jobs: %v", err)
			SendMessage(bot, message.Chat.ID, "❌ Sorry, I couldn't fetch jobs right now. Please try again later.")
			return
		}

		if len(jobResp.Jobs) == 0 {
			SendMessage(bot, message.Chat.ID, "😔 No jobs found with the specified criteria. Try different filters.")
			return
		}
		// Send jobs (split into multiple messages if needed)
		var messageBuffer strings.Builder
		messageBuffer.WriteString(fmt.Sprintf("🎯 Found %d remote jobs:\n\n", len(jobResp.Jobs)))

		for i, job := range jobResp.Jobs {
			jobMsg := services.FormatJobMsg(job)

			// Check if adding this job would exceed Telegram's message limit (4096 characters)
			if messageBuffer.Len()+len(jobMsg) > 4000 {
				// Send current buffer
				SendMessage(bot, message.Chat.ID, messageBuffer.String())

				// Start new buffer
				messageBuffer.Reset()
				messageBuffer.WriteString(fmt.Sprintf("📋 Continuing jobs list (%d/%d):\n\n", i+1, len(jobResp.Jobs)))
			}

			messageBuffer.WriteString(jobMsg)
		}
		// Send remaining jobs
		if messageBuffer.Len() > 0 {
			SendMessage(bot, message.Chat.ID, messageBuffer.String())
		}

	default:
		SendMessage(bot, message.Chat.ID, "❓ Unknown command. Type /help for available commands.")

	}
}
