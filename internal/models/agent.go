package models

import "time"

type Agent struct {
	ID                string      `json:"id"`
	ProjectID         *string     `json:"project_id"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	Status            AgentStatus `json:"status"`
	Source            AgentSource `json:"source"`
	SystemPrompt      string      `json:"system_prompt"`
	WorkDir           string      `json:"work_dir"`
	ExitCode          *int        `json:"exit_code"`
	ErrorMessage      *string     `json:"error_message"`
	GitHubIssueNumber *int        `json:"github_issue_number"`
	GitHubPRNumber    *int        `json:"github_pr_number"`
	PID               *int        `json:"pid"`
	StartedAt         *time.Time  `json:"started_at"`
	FinishedAt        *time.Time  `json:"finished_at"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type AgentLog struct {
	ID        int       `json:"id"`
	AgentID   string    `json:"agent_id"`
	Stream    LogStream `json:"stream"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
