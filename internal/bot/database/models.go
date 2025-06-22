package database

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

// Jobicy API structs
type JobicyResponse struct { // Job is the API response, Count the number of listings from query params
	Jobs  []Job `json:"jobs"`
	Count int   `json:"count"`
	// add if to customize
}

// Response
type Job struct {
	ID             int    `json:"id"`
	URL            string `json:"url"`
	JobTitle       string `json:"jobTitle"`
	CompanyName    string `json:"companyName"`
	CompanyLogo    string `json:"companyLogo"`
	JobType        string `json:"jobType"`
	JobGeo         string `json:"jobGeo"`
	JobLevel       string `json:"jobLevel"`
	JobExcerpt     string `json:"jobExcerpt"`
	JobDescription string `json:"jobDescription"`
	PubDate        string `json:"pubDate"`
	Salary_min     string `json:"salaryMin,omitempty"`
	Salary_max     string `json:"salaryMax,omitempty"`
	// Add more fields as needed from Jobicy API
}

// TG Conn. Pool
type TelegramBot struct {
	Token   string
	BaseURL string
	Offset  int
}
