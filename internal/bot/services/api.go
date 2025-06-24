package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tg-remote/internal/types"
)

// NewJobicyClient creates a new Jobicy client
func NewJobicyClient() *types.Jobicy {
	return &types.Jobicy{
		BaseURL: "http://jobicy.com/api/v2/remote-jobs",
	}
}

// GetJobs fetches jobs from Jobicy API
func GetJobs(client *types.Jobicy, count int, geo string, industry string, tag string) (*types.JobicyResponse, error) {
	url := fmt.Sprintf("%s?count=%d", client.BaseURL, count)

	// add optional params
	// since it is for educational purposes, and some params like geo
	// doesn't apply to certain geolocations, lets leave it at that
	// afterall, they are optional
	if geo != "" {
		url += "&geo=" + geo
	}
	if industry != "" {
		url += "&industry=" + industry
	}
	if tag != "" {
		url += "&tag=" + tag
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

	var jobicyResp types.JobicyResponse
	if err := json.Unmarshal(body, &jobicyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &jobicyResp, nil
}
