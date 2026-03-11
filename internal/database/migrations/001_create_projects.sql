CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    root_folder TEXT NOT NULL,
    allowed_emails TEXT NOT NULL DEFAULT '[]',
    github_owner TEXT,
    github_repo TEXT,
    github_webhook_secret TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_projects_github ON projects(github_owner, github_repo) WHERE github_owner IS NOT NULL AND github_repo IS NOT NULL;
