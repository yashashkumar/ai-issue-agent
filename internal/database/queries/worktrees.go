package queries

const (
	CreateWorktree = `
		INSERT INTO worktrees (
			id, project_id, agent_id, branch_name, worktree_path,
			github_issue_number, is_active, created_at, cleaned_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	GetActiveWorktreeForIssue = `
		SELECT
			id, project_id, agent_id, branch_name, worktree_path,
			github_issue_number, is_active, created_at, cleaned_at
		FROM worktrees
		WHERE project_id = ? AND github_issue_number = ? AND is_active = 1
	`

	MarkWorktreeCleaned = `
		UPDATE worktrees SET
			is_active = 0, cleaned_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
)
