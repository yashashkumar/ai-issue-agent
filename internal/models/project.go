package models

import (
	"encoding/json"
	"time"
)

type Project struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	RootFolder          string    `json:"root_folder"`
	AllowedEmails       []string  `json:"allowed_emails"` // Stored as JSON string in DB
	GitHubOwner         *string   `json:"github_owner"`
	GitHubRepo          *string   `json:"github_repo"`
	GitHubWebhookSecret *string   `json:"-"` // Omit from generic JSON output for security
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// SetAllowedEmails helper serializes slice to JSON byte string
func (p *Project) AllowedEmailsJSON() string {
	b, _ := json.Marshal(p.AllowedEmails)
	return string(b)
}

// ParseAllowedEmails helper deserializes JSON string to slice
func ParseAllowedEmails(jsonStr string) []string {
	var emails []string
	if err := json.Unmarshal([]byte(jsonStr), &emails); err != nil {
		return []string{}
	}
	return emails
}
