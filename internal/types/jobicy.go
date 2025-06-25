package types

// Jobicy API client
type Jobicy struct {
	BaseURL string
}

// JobicyResponse represents the API response structure
// should be data absolutely necessary to display!
// sometimes the structure has its data contained in an array
type JobicyResponse struct {
	Jobs []Job `json:"jobs"`
}

// Job represents a single job posting out of the response
type Job struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	JobTitle    string `json:"jobTitle"`
	CompanyName string `json:"companyName"`
	JobGeo      string `json:"jobGeo"`
	//JobType     string `json:"jobType"`	// supposed to be an array of type of jobs
	// add job industry
	JobLevel   string `json:"jobLevel"`
	JobExcerpt string `json:"jobExcerpt"` // replace with jobdescription and truncate excessively long ones?
	PubDate    string `json:"pubDate"`
	Salary_min string `json:"salary_min"`
	Salary_max string `json:"salary_max"`
	// salary currency
	// salaryperiod
}
