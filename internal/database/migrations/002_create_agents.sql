CREATE TABLE agents (
    id TEXT PRIMARY KEY,
    project_id TEXT,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'running', 'finished', 'error')),
    source TEXT NOT NULL CHECK(source IN ('github_issue', 'github_comment', 'custom')),
    system_prompt TEXT NOT NULL,
    work_dir TEXT NOT NULL,
    exit_code INTEGER,
    error_message TEXT,
    github_issue_number INTEGER,
    github_pr_number INTEGER,
    pid INTEGER,
    started_at DATETIME,
    finished_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE SET NULL
);

CREATE INDEX idx_agents_project_id ON agents(project_id);
CREATE INDEX idx_agents_status ON agents(status);
CREATE INDEX idx_agents_source ON agents(source);
CREATE INDEX idx_agents_github_issue ON agents(github_issue_number);
