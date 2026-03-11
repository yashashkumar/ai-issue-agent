package queries

const (
	CreateProject = `
		INSERT INTO projects (
			id, name, description, root_folder, allowed_emails,
			github_owner, github_repo, github_webhook_secret,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	UpdateProject = `
		UPDATE projects SET
			name = ?, description = ?, root_folder = ?, allowed_emails = ?,
			github_owner = ?, github_repo = ?, github_webhook_secret = ?,
			updated_at = ?
		WHERE id = ?
	`

	DeleteProject = `DELETE FROM projects WHERE id = ?`

	GetProjectByID = `
		SELECT
			id, name, description, root_folder, allowed_emails,
			github_owner, github_repo, github_webhook_secret,
			created_at, updated_at
		FROM projects
		WHERE id = ?
	`

	GetProjectByGitHubRepo = `
		SELECT
			id, name, description, root_folder, allowed_emails,
			github_owner, github_repo, github_webhook_secret,
			created_at, updated_at
		FROM projects
		WHERE github_owner = ? AND github_repo = ?
	`

	ListProjects = `
		SELECT
			id, name, description, root_folder, allowed_emails,
			github_owner, github_repo, github_webhook_secret,
			created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`

	SearchProjects = `
		SELECT
			id, name, description, root_folder, allowed_emails,
			github_owner, github_repo, github_webhook_secret,
			created_at, updated_at
		FROM projects
		WHERE name LIKE ? OR description LIKE ?
		ORDER BY created_at DESC
	`
)
