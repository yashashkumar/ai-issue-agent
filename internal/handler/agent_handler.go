package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/yashashkumar/ai-issue-agent/internal/service"
	"github.com/yashashkumar/ai-issue-agent/internal/spawner"
)

type AgentHandler struct {
	agentSvc   *service.AgentService
	projectSvc *service.ProjectService
	spawner    *spawner.Spawner
}

func NewAgentHandler(agentSvc *service.AgentService, projectSvc *service.ProjectService, spawner *spawner.Spawner) *AgentHandler {
	return &AgentHandler{
		agentSvc:   agentSvc,
		projectSvc: projectSvc,
		spawner:    spawner,
	}
}

func (h *AgentHandler) ListProjectAgents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)
	projectID := r.PathValue("id")

	agents, err := h.agentSvc.ListProjectAgents(ctx, projectID)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "failed to fetch project agents", reqID)
		return
	}

	WriteJSON(w, http.StatusOK, agents)
}

func (h *AgentHandler) ListAllAgents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)

	agents, err := h.agentSvc.ListAllAgents(ctx)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "failed to fetch agents", reqID)
		return
	}

	WriteJSON(w, http.StatusOK, agents)
}

func (h *AgentHandler) GetAgent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)
	agentID := r.PathValue("id")

	agent, err := h.agentSvc.GetAgent(ctx, agentID)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, "agent not found", reqID)
		return
	}

	WriteJSON(w, http.StatusOK, agent)
}

func (h *AgentHandler) GetAgentLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)
	agentID := r.PathValue("id")
	stream := r.URL.Query().Get("stream")

	logs, err := h.agentSvc.GetAgentLogs(ctx, agentID, stream)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "failed to fetch logs", reqID)
		return
	}

	WriteJSON(w, http.StatusOK, logs)
}

func (h *AgentHandler) CancelAgent(w http.ResponseWriter, r *http.Request) {
	// The cancel triggers the spawner to terminate
	reqID := GetRequestID(r.Context())
	agentID := r.PathValue("id")

	if err := h.spawner.CancelAgent(agentID); err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "failed to cancel agent", reqID)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// Minimal implementation for SSE endpoint if needed
func (h *AgentHandler) StreamAgentLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	agentID := r.PathValue("id")
	// Phase 5 will expand SSE here by reading from a pub/sub mechanism in Spawner
	// Mock stream for now:
	for i := 0; i < 5; i++ {
		fmt.Fprintf(w, "data: %s - streaming %s\n\n", time.Now().Format(time.RFC3339), agentID)
		flusher.Flush()
		select {
		case <-r.Context().Done():
			return
		case <-time.After(1 * time.Second):
		}
	}
}

// Endpoints for files
func (h *AgentHandler) ListAgentFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)
	agentID := r.PathValue("id")

	agent, err := h.agentSvc.GetAgent(ctx, agentID)
	if err != nil || agent == nil {
		WriteJSONError(w, http.StatusNotFound, "agent not found", reqID)
		return
	}

	// Recursively list files in agent.WorkDir
	type FileInfo struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Size int64  `json:"size"`
	}

	var files []FileInfo
	filepath.Walk(agent.WorkDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(agent.WorkDir, path)
		files = append(files, FileInfo{
			Name: info.Name(),
			Path: relPath,
			Size: info.Size(),
		})
		return nil
	})

	WriteJSON(w, http.StatusOK, files)
}
