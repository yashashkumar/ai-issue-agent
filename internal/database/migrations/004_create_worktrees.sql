CREATE TABLE worktrees (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    agent_id TEXT NOT NULL,
    branch_name TEXT NOT NULL,
    worktree_path TEXT NOT NULL,
    github_issue_number INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cleaned_at DATETIME,
    FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY(agent_id) REFERENCES agents(id)
);

CREATE UNIQUE INDEX idx_worktrees_active_issue ON worktrees(project_id, github_issue_number) WHERE is_active = 1;
