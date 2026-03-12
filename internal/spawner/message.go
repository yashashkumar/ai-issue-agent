package spawner

import "github.com/yashashkumar/ai-issue-agent/internal/models"

type SpawnRequest struct {
	ID                string
	Source            models.AgentSource
	ProjectID         *string
	SystemPrompt      string
	WorkDir           string
	Name              string
	Description       string
	GitHubIssueNumber *int
	GitHubOwner       *string
	GitHubRepo        *string
	CreateBranch      bool
	BranchName        *string
	ExistingWorktree  *string
	CreatePR          bool
}
