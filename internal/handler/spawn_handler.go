package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/yashashkumar/ai-issue-agent/internal/models"
	"github.com/yashashkumar/ai-issue-agent/internal/service"
	"github.com/yashashkumar/ai-issue-agent/internal/spawner"
)

type SpawnHandler struct {
	spawner     *spawner.Spawner
	projectSvc  *service.ProjectService
	workBaseDir string
}

func NewSpawnHandler(spawner *spawner.Spawner, projectSvc *service.ProjectService, workBaseDir string) *SpawnHandler {
	return &SpawnHandler{
		spawner:     spawner,
		projectSvc:  projectSvc,
		workBaseDir: workBaseDir,
	}
}

type SpawnRequestPayload struct {
	SystemPrompt string  `json:"system_prompt"`
	ProjectID    *string `json:"project_id"`
	WorkDir      *string `json:"work_dir"`
	Name         *string `json:"name"`
	Description  *string `json:"description"`
}

func (h *SpawnHandler) HandleSpawn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)

	var payload SpawnRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid request body", reqID)
		return
	}

	if payload.SystemPrompt == "" || len(payload.SystemPrompt) < 10 {
		WriteJSONError(w, http.StatusBadRequest, "system_prompt is required and must be at least 10 chars", reqID)
		return
	}

	var resolvedWorkDir string
	if payload.WorkDir != nil && *payload.WorkDir != "" {
		resolvedWorkDir = *payload.WorkDir
		if info, err := os.Stat(resolvedWorkDir); err != nil || !info.IsDir() {
			WriteJSONError(w, http.StatusBadRequest, "work_dir does not exist or is not a directory", reqID)
			return
		}
	} else if payload.ProjectID != nil && *payload.ProjectID != "" {
		proj, err := h.projectSvc.GetProject(ctx, *payload.ProjectID)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, "invalid project_id", reqID)
			return
		}
		resolvedWorkDir = proj.RootFolder
	} else {
		// Use temp directory under base
		resolvedWorkDir = h.workBaseDir + "/temp_" + uuid.NewString()[:8]
		os.MkdirAll(resolvedWorkDir, 0755) // Ignore errors, spawner will handle it
	}

	agentID := uuid.NewString()
	name := "Custom Agent " + agentID[:8]
	if payload.Name != nil && *payload.Name != "" {
		name = *payload.Name
	}
	desc := ""
	if payload.Description != nil {
		desc = *payload.Description
	}

	req := spawner.SpawnRequest{
		ID:           agentID,
		Source:       models.AgentSourceCustom,
		ProjectID:    payload.ProjectID,
		SystemPrompt: payload.SystemPrompt,
		WorkDir:      resolvedWorkDir,
		Name:         name,
		Description:  desc,
	}

	// Dispatch to spawner
	h.spawner.Spawn(req)

	response := map[string]string{
		"agent_id": agentID,
		"status":   "pending",
	}
	WriteJSON(w, http.StatusAccepted, response)
}
