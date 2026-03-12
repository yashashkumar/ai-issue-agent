package models

type AgentStatus string

const (
	AgentStatusPending  AgentStatus = "pending"
	AgentStatusRunning  AgentStatus = "running"
	AgentStatusFinished AgentStatus = "finished"
	AgentStatusError    AgentStatus = "error"
)

type AgentSource string

const (
	AgentSourceGitHubIssue   AgentSource = "github_issue"
	AgentSourceGitHubComment AgentSource = "github_comment"
	AgentSourceCustom        AgentSource = "custom"
)

type LogStream string

const (
	LogStreamStdout LogStream = "stdout"
	LogStreamStderr LogStream = "stderr"
)
