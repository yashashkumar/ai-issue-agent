package queries

const (
	CreateAgent = `
		INSERT INTO agents (
			id, project_id, name, description, status, source,
			system_prompt, work_dir, exit_code, error_message,
			github_issue_number, github_pr_number, pid,
			started_at, finished_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	UpdateAgentStatus = `
		UPDATE agents SET
			status = ?, error_message = ?, exit_code = ?,
			started_at = COALESCE(?, started_at), finished_at = COALESCE(?, finished_at),
			updated_at = ?
		WHERE id = ?
	`

	UpdateAgentPID = `
		UPDATE agents SET pid = ?, updated_at = ? WHERE id = ?
	`

	UpdateAgentPRNumber = `
		UPDATE agents SET github_pr_number = ?, updated_at = ? WHERE id = ?
	`

	GetAgentByID = `
		SELECT
			id, project_id, name, description, status, source,
			system_prompt, work_dir, exit_code, error_message,
			github_issue_number, github_pr_number, pid,
			started_at, finished_at, created_at, updated_at
		FROM agents
		WHERE id = ?
	`

	ListProjectAgents = `
		SELECT
			id, project_id, name, description, status, source,
			system_prompt, work_dir, exit_code, error_message,
			github_issue_number, github_pr_number, pid,
			started_at, finished_at, created_at, updated_at
		FROM agents
		WHERE project_id = ?
		ORDER BY created_at DESC
	`

	ListAllAgents = `
		SELECT
			id, project_id, name, description, status, source,
			system_prompt, work_dir, exit_code, error_message,
			github_issue_number, github_pr_number, pid,
			started_at, finished_at, created_at, updated_at
		FROM agents
		ORDER BY created_at DESC
	`

	InsertAgentLog = `
		INSERT INTO agent_logs (agent_id, stream, content, created_at)
		VALUES (?, ?, ?, ?)
	`

	GetAgentLogs = `
		SELECT id, agent_id, stream, content, created_at
		FROM agent_logs
		WHERE agent_id = ? AND (stream = ? OR ? = 'all')
		ORDER BY created_at ASC
	`
)
