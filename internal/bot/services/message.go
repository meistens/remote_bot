package services

import (
	"fmt"
	"strings"
	"tg-remote/internal/bot"
)

// Format for TG msg, generated skello
func formatJobMsg(job bot.Job) string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("<b>* %s</b>\n", job.JobTitle))

	message.WriteString(fmt.Sprintf("<b>*Company:</b> %s\n", job.CompanyName))

	// TODO: Add your desired fields here
	// Examples of fields you can use:
	message.WriteString(fmt.Sprintf("<b>📍 Location:</b> %s\n", job.JobGeo))
	message.WriteString(fmt.Sprintf("<b>💼 Type:</b> %s\n", job.JobType))
	message.WriteString(fmt.Sprintf("<b>📊 Level:</b> %s\n", job.JobLevel))

	// Salary fields (if available)
	if job.Salary_min != "" || job.Salary_max != "" {
		message.WriteString(fmt.Sprintf("<b>💰 Salary:</b> %s - %s\n",
			job.Salary_min, job.Salary_max))
	}

	// Job excerpt/description
	if job.JobExcerpt != "" {
		excerpt := job.JobExcerpt
		if len(excerpt) > 200 {
			excerpt = excerpt[:200] + "..."
		}
		message.WriteString(fmt.Sprintf("<b>📝 Description:</b> %s\n", excerpt))
	}

	// Publication date
	message.WriteString(fmt.Sprintf("<b>📅 Posted:</b> %s\n", job.PubDate))

	// Job URL
	message.WriteString(fmt.Sprintf("<b>🔗 Apply:</b> <a href=\"%s\">View Job</a>\n", job.URL))

	// Add separator
	message.WriteString("\n" + strings.Repeat("─", 30) + "\n\n")

	return message.String()

}
