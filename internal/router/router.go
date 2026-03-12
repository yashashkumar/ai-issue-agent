package router

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/yashashkumar/ai-issue-agent/internal/handler"
)

type RouterDependencies struct {
	ProjectHandler *handler.ProjectHandler
	AgentHandler   *handler.AgentHandler
	SpawnHandler   *handler.SpawnHandler
	WebhookHandler *handler.WebhookHandler
	Logger         *slog.Logger
}

func NewRouter(deps RouterDependencies) *http.ServeMux {
	mux := http.NewServeMux()

	// This is a helper to chain our generic middlewares
	chain := func(h http.HandlerFunc) http.Handler {
		var next http.Handler = h
		next = handler.CORSMiddleware(next)
		next = handler.LoggingMiddleware(deps.Logger)(next)
		next = handler.RecoveryMiddleware(deps.Logger)(next)
		next = handler.RequestIDMiddleware(next)
		return next
	}

	// Health Check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		handler.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Gateway Webhook endpoints
	mux.Handle("POST /gateway/hooks/spawn/gh/{owner}/{repo}/issues", chain(deps.WebhookHandler.HandleGitHubIssues))
	mux.Handle("POST /gateway/hooks/spawn/gh/{owner}/{repo}/issue-comments", chain(deps.WebhookHandler.HandleGitHubIssueComments))

	// API - Custom Spawn
	mux.Handle("POST /api/v1/spawn", chain(deps.SpawnHandler.HandleSpawn))

	// API - Projects
	mux.Handle("GET /api/v1/projects", chain(deps.ProjectHandler.ListProjects))
	mux.Handle("POST /api/v1/projects", chain(deps.ProjectHandler.CreateProject))
	mux.Handle("GET /api/v1/projects/{id}", chain(deps.ProjectHandler.GetProject))
	mux.Handle("PUT /api/v1/projects/{id}", chain(deps.ProjectHandler.UpdateProject))
	mux.Handle("DELETE /api/v1/projects/{id}", chain(deps.ProjectHandler.DeleteProject))

	// API - Agents
	mux.Handle("GET /api/v1/projects/{id}/agents", chain(deps.AgentHandler.ListProjectAgents))
	mux.Handle("GET /api/v1/agents", chain(deps.AgentHandler.ListAllAgents))
	mux.Handle("GET /api/v1/agents/{id}", chain(deps.AgentHandler.GetAgent))
	mux.Handle("GET /api/v1/agents/{id}/logs", chain(deps.AgentHandler.GetAgentLogs))
	mux.Handle("GET /api/v1/agents/{id}/logs/stream", chain(deps.AgentHandler.StreamAgentLogs))
	mux.Handle("POST /api/v1/agents/{id}/cancel", chain(deps.AgentHandler.CancelAgent))
	mux.Handle("GET /api/v1/agents/{id}/files", chain(deps.AgentHandler.ListAgentFiles))

	// UI - SPA Endpoints
	spa := &spaHandler{staticPath: "frontend/dist", indexPath: "index.html"}
	mux.Handle("/", chain(spa.ServeHTTP))

	return mux
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Check if the file exists in the static dir
	_, err := http.Dir(h.staticPath).Open(path)
	if os.IsNotExist(err) {
		// File does not exist, serve index.html
		http.ServeFile(w, r, h.staticPath+"/"+h.indexPath)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// File exists, serve it
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
