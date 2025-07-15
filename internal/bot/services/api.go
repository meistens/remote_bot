package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"tg-remote/internal/types"
)

// NewJobicyClient creates a new Jobicy client
func NewJobicyClient() *types.Jobicy {
	return &types.Jobicy{
		BaseURL: "https://jobicy.p.rapidapi.com/api/v2/remote-jobs",
	}
}

// GetJobs fetches jobs from Jobicy API via RapidAPI
func GetJobs(client *types.Jobicy, count int, geo string, industry string, tag string) (*types.JobicyResponse, error) {
	url := fmt.Sprintf("%s?count=%d", client.BaseURL, count)

	// Add optional params
	if geo != "" {
		url += "&geo=" + geo
	}
	if industry != "" {
		url += "&industry=" + industry
	}
	if tag != "" {
		url += "&tag=" + tag
	}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add RapidAPI headers
	rapidAPIKey := os.Getenv("RAPIDAPI_KEY")
	if rapidAPIKey == "" {
		return nil, fmt.Errorf("RAPIDAPI_KEY environment variable is required")
	}

	req.Header.Set("X-RapidAPI-Key", rapidAPIKey)
	req.Header.Set("X-RapidAPI-Host", "jobicy.p.rapidapi.com")

	// Make request
	client_http := &http.Client{}
	resp, err := client_http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var jobicyResp types.JobicyResponse
	if err := json.Unmarshal(body, &jobicyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &jobicyResp, nil
}
