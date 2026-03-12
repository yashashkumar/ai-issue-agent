export type Project = {
  id: string;
  name: string;
  description: string;
  root_folder: string;
  allowed_emails: string[];
  github_owner?: string | null;
  github_repo?: string | null;
  created_at: string;
  updated_at: string;
};

export type AgentStatus = 'pending' | 'running' | 'finished' | 'error';
export type AgentSource = 'github_issue' | 'github_comment' | 'custom';

export type Agent = {
  id: string;
  project_id?: string | null;
  name: string;
  description: string;
  status: AgentStatus;
  source: AgentSource;
  system_prompt: string;
  work_dir: string;
  exit_code?: number | null;
  error_message?: string | null;
  github_issue_number?: number | null;
  github_pr_number?: number | null;
  pid?: number | null;
  started_at?: string | null;
  finished_at?: string | null;
  created_at: string;
  updated_at: string;
};

export type AgentLog = {
  id: number;
  agent_id: string;
  stream: 'stdout' | 'stderr';
  content: string;
  created_at: string;
};
