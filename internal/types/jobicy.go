package types

// Jobicy API client
type Jobicy struct {
	BaseURL string
}

// JobicyResponse represents the API response structure
type JobicyResponse struct {
	Jobs []Job `json:"jobs"`
}

// Job represents a single job posting
type Job struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	JobTitle    string `json:"jobTitle"`
	CompanyName string `json:"companyName"`
	JobGeo      string `json:"jobGeo"`
	//JobType     string `json:"jobType"`	// supposed to be an array of type of jobs
	JobLevel   string `json:"jobLevel"`
	JobExcerpt string `json:"jobExcerpt"`
	PubDate    string `json:"pubDate"`
	Salary_min string `json:"salary_min"`
	Salary_max string `json:"salary_max"`
}
