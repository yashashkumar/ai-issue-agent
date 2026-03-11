package models

import "time"

type Worktree struct {
	ID                string     `json:"id"`
	ProjectID         string     `json:"project_id"`
	AgentID           string     `json:"agent_id"`
	BranchName        string     `json:"branch_name"`
	WorktreePath      string     `json:"worktree_path"`
	GitHubIssueNumber int        `json:"github_issue_number"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
	CleanedAt         *time.Time `json:"cleaned_at"`
}
