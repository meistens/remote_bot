package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tg-remote/internal/bot"
)

// Jobicy API client
type Jobicy struct {
	BaseURL string
}

// Jobicy
func NewJobicyClient() *Jobicy {
	return &Jobicy{
		BaseURL: "http://jobicy.com/api/v2/remote-jobs",
	}
}

func (client *Jobicy) GetJobs(count int, geo string, industry string, tag string) (*bot.JobicyResponse, error) {
	url := fmt.Sprintf("%s?count=%d", client.BaseURL, count)

	// add optional params, incase for tweaking
	if geo != "" {
		url += "&geo=" + geo
	}
	if industry != "" {
		url += "&industry" + industry
	}
	if tag != "" {
		url += "&tag" + tag
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var jobicyResp bot.JobicyResponse
	if err := json.Unmarshal(body, &jobicyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &jobicyResp, nil
}
