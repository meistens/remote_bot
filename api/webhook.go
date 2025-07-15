package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// TODO: remove other files
// Telegram Bot structures
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
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

type SendMessagePayload struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// Jobicy API structures
type JobicyResponse struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	JobTitle    string `json:"jobTitle"`
	CompanyName string `json:"companyName"`
	JobGeo      string `json:"jobGeo"`
	JobLevel    string `json:"jobLevel"`
	JobExcerpt  string `json:"jobExcerpt"`
	PubDate     string `json:"pubDate"`
	Salary_min  string `json:"salary_min"`
	Salary_max  string `json:"salary_max"`
}

// TelegramBot client
type TelegramBot struct {
	Token   string
	BaseURL string
}

func NewTelegramBot(token string) *TelegramBot {
	return &TelegramBot{
		Token:   token,
		BaseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
	}
}

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
	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
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

// Jobicy API client
func GetJobs(count int, geo string, industry string, tag string) (*JobicyResponse, error) {
	url := fmt.Sprintf("https://jobicy.p.rapidapi.com/api/v2/remote-jobs?count=%d", count)

	if geo != "" {
		url += "&geo=" + geo
	}
	if industry != "" {
		url += "&industry=" + industry
	}
	if tag != "" {
		url += "&tag=" + tag
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	rapidAPIKey := os.Getenv("RAPIDAPI_KEY")
	if rapidAPIKey == "" {
		return nil, fmt.Errorf("RAPIDAPI_KEY environment variable is required")
	}

	req.Header.Set("X-RapidAPI-Key", rapidAPIKey)
	req.Header.Set("X-RapidAPI-Host", "jobicy.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var jobicyResp JobicyResponse
	if err := json.Unmarshal(body, &jobicyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &jobicyResp, nil
}

// Message formatting
func FormatJobMsg(job Job) string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("<b>%s</b>\n", job.JobTitle))
	message.WriteString(fmt.Sprintf("<b>Company:</b> %s\n", job.CompanyName))
	message.WriteString(fmt.Sprintf("<b>Location:</b> %s\n", job.JobGeo))
	message.WriteString(fmt.Sprintf("<b>Level:</b> %s\n", job.JobLevel))

	if job.Salary_min != "" || job.Salary_max != "" {
		message.WriteString(fmt.Sprintf("<b>Salary:</b> %s - %s\n", job.Salary_min, job.Salary_max))
	}

	if job.JobExcerpt != "" {
		excerpt := job.JobExcerpt
		if len(excerpt) > 200 {
			excerpt = excerpt[:200] + "..."
		}
		message.WriteString(fmt.Sprintf("<b>Description:</b> %s\n", excerpt))
	}

	message.WriteString(fmt.Sprintf("<b>Posted On:</b> %s\n", strings.Split(job.PubDate, "T")[0]))
	message.WriteString(fmt.Sprintf("<b>Apply:</b> <a href=\"%s\">View Job</a>\n", job.URL))
	message.WriteString("\n" + strings.Repeat("─", 30) + "\n\n")

	return message.String()
}

// Command parsing
func ParseCommand(text string) (command string, args map[string]string) {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return "", nil
	}

	command = parts[0]
	args = make(map[string]string)

	for i := 1; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "--") {
			key := strings.TrimPrefix(parts[i], "--")
			if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "--") {
				args[key] = parts[i+1]
				i++
			} else {
				args[key] = "true"
			}
		}
	}
	return command, args
}

// Handler is the main webhook handler for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// html
	w.Header().Set("Content-Type", "text/html")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Allow only POST
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

	// Parse webhook update
	var update Update
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

	// Create bot instance
	telegramBot := NewTelegramBot(botToken)

	// If a message exists
	if update.Message.Text != "" {
		log.Printf("Received message from %s: %s", update.Message.From.Username, update.Message.Text)

		// Handle message
		err := handleMessage(telegramBot, update.Message)
		if err != nil {
			log.Printf("Error handling message: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Remote Jobs Telegram Bot</title>
			<style>
				body { font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; }
				.status { color: green; font-weight: bold; }
				.endpoint { background: #f5f5f5; padding: 10px; border-radius: 5px; margin: 10px 0; }
				code { background: #f0f0f0; padding: 2px 4px; border-radius: 3px; }
			</style>
		</head>
		<body>
			<h1>Remote Jobs Telegram Bot</h1>
			<p class="status">✅ Bot is running and ready to receive webhooks!</p>

			<h2>How to Use:</h2>
			<ol>
				<li>Find the bot on Telegram</li>
				<li>Send <code>/start</code> to begin</li>
				<li>Use <code>/jobs</code> to search for remote jobs</li>
				<li>Use <code>/help</code> for more commands</li>
			</ol>

			<h2>Available Commands:</h2>
			<ul>
				<li><code>/jobs</code> - Get latest remote jobs</li>
				<li><code>/jobs --count 10</code> - Get specific number of jobs</li>
				<li><code>/jobs --geo USA</code> - Filter by location</li>
				<li><code>/jobs --tag python</code> - Filter by technology</li>
			</ul>

			<div class="endpoint">
				<strong>Webhook Endpoint:</strong> <code>POST /api/webhook</code>
			</div>

			<p><em>This page confirms your bot is deployed successfully.</em></p>
		</body>
		</html>
		`

	fmt.Fprint(w, html)

	// Return status
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleMessage(bot *TelegramBot, message Message) error {
	command, args := ParseCommand(message.Text)

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

func handleJobsCommand(bot *TelegramBot, message Message, args map[string]string) error {
	// Parse args
	count := 5
	if countStr, ok := args["count"]; ok {
		if c, err := strconv.Atoi(countStr); err == nil && c > 0 && c <= 50 {
			count = c
		}
	}

	geo := args["geo"]
	industry := args["industry"]
	tag := args["tag"]

	// Send searching message
	err := bot.SendMessage(message.Chat.ID, "🔍 Searching for jobs...")
	if err != nil {
		return err
	}

	// Fetch jobs from API
	jobResp, err := GetJobs(count, geo, industry, tag)
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
		jobMsg := FormatJobMsg(job)

		// Check if adding this job would exceed Telegram's message limit
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
