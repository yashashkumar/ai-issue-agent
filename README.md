# AgentForge

AgentForge is an HTTP gateway and orchestration platform that spawns and manages `gemini` CLI AI agents. Built with Go 1.23+, pure stdlib HTTP handling, SQLite, and a React + Vite frontend.

## Features
- **Gateway**: Receive GitHub Webhooks (issues/issue comments) to spawn AI agents automatically.
- **Process Orchestration**: Limit max concurrency, capture logs (stdout/stderr) natively from CLI processes, isolate processes and cleanly shut down.
- **Git Worktree Manager**: Link GitHub issue threads to specific branches, so agents hold contextual worktrees without rewriting the repo from scratch.
- **Live Log Streaming**: Watch agent output in real time from the AgentDetail page.
- **React Dashboard UI**: Modern SPA to spawn, manage, and audit all AI tasks.

---

## Prerequisites

Before running AgentForge, ensure you have the following installed:

| Dependency | Version | Notes |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.23+ | Backend runtime |
| [Node.js](https://nodejs.org/) | **v20+ required** (v24 recommended) | For frontend build |
| [nvm](https://github.com/nvm-sh/nvm) | any | Recommended for managing Node versions |
| Gemini CLI | latest | `npm install -g @google/gemini-cli` |

> **Important:** The `gemini` CLI requires **Node.js v20 or later**. It will crash with a `SyntaxError: Invalid regular expression flags` on Node v18 or earlier.

---

## First-Time Setup

### 1. Switch to Node v24

```bash
source ~/.nvm/nvm.sh
nvm use 24.1.0        # or: nvm install 24 && nvm use 24
```

### 2. Install the Gemini CLI on Node v24

```bash
npm install -g @google/gemini-cli
gemini --version      # should print 0.1.x or later
```

### 3. Authenticate with Gemini (one-time)

Run the CLI interactively **once** to complete the Google OAuth login. This stores credentials in `~/.gemini/oauth_creds.json`.

```bash
gemini
# Choose "Login with Google" and complete the browser flow
# Then type /quit to exit
```

> **This step is mandatory.** AgentForge spawns `gemini` as a subprocess and reads the stored OAuth token. If you skip this, agents will fail with "Please set an Auth method".

### 4. Configure environment

Copy the example env file and adjust as needed:

```bash
cp .env.example .env
```

Key variables in `.env`:

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DATABASE_PATH=./data/agentforge.db
MAX_CONCURRENT_AGENTS=5
AGENT_TIMEOUT_MINUTES=30
DEFAULT_WORK_BASE_DIR=./workspaces
GEMINI_CLI_PATH=gemini
```

If you prefer API key auth over Google login, also set:
```env
GEMINI_API_KEY=AIza...your_key_here
```

---

## Running the App

> **Important:** Always run `make run` from the **same terminal session** where you ran `gemini` (so Node v24's nvm environment is active). This ensures the Go backend inherits the correct PATH and OAuth credentials.

```bash
# Kill any existing process on port 8080 (if needed)
npx kill-port 8080

# Build the React frontend + Go backend, then start the server
make run
```

Open [http://localhost:8080](http://localhost:8080) in your browser.

---

## Makefile Targets

| Target | Description |
|---|---|
| `make build` | Build the React frontend (`npm run build`) and Go binary |
| `make run` | Build everything and start the server |
| `make clean` | Remove binary, workspace dirs, frontend dist, and database |
| `make test` | Run Go tests |

---

## Troubleshooting

### `SyntaxError: Invalid regular expression flags`
**Cause:** `gemini` is running under Node v18 or lower.  
**Fix:** Switch to Node v20+ via nvm: `nvm use 24` then reinstall gemini: `npm install -g @google/gemini-cli`

### `Please set an Auth method in ~/.gemini/settings.json`
**Cause:** Either the OAuth token is missing, or `make run` was started in a terminal without the Node v24 nvm environment active.  
**Fix:**
1. Run `gemini` interactively in the **same terminal** first (complete Google login).
2. Then run `make run` from that **same terminal** session.

### `listen tcp 0.0.0.0:8080: bind: address already in use`
**Cause:** A previous server instance is still running.  
**Fix:** `npx kill-port 8080` then `make run` again.

---

## Architecture Overview

```
frontend/           React + Vite SPA (built to frontend/dist/)
internal/
  router/           http.NewServeMux routes + SPA fallback handler
  handler/          HTTP controller structs
  service/          Business logic (AgentService, ProjectService)
  repository/       Pure database/sql queries (no ORM)
  spawner/          Concurrency manager — channels, semaphores, goroutines
  database/         SQLite via modernc.org/sqlite + embedded migrations
cmd/server/         main.go entrypoint
```

---

## API Reference

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/spawn` | Trigger a standalone agent spawn with a custom prompt |
| `POST` | `/api/v1/projects` | Create a new project |
| `GET` | `/api/v1/projects` | List all projects |
| `GET` | `/api/v1/projects/{id}` | Get project details |
| `DELETE` | `/api/v1/projects/{id}` | Delete a project |
| `GET` | `/api/v1/agents` | List all agents |
| `GET` | `/api/v1/agents/{id}` | Get agent details |
| `GET` | `/api/v1/agents/{id}/logs` | Fetch stored agent log output |
| `GET` | `/api/v1/agents/{id}/logs/stream` | Stream agent logs (SSE) |
| `POST` | `/api/v1/agents/{id}/cancel` | Send SIGKILL to a running agent |
| `GET` | `/api/v1/agents/{id}/files` | List files in the agent's workspace |

### Webhooks

```
POST /gateway/hooks/spawn/gh/{owner}/{repo}/issues
POST /gateway/hooks/spawn/gh/{owner}/{repo}/issue-comments
```

Requires `X-Hub-Signature-256` header. Listens for `opened`/`edited` issue events and `created`/`edited` comment events.
