package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/yashashkumar/ai-issue-agent/internal/models"
	"github.com/yashashkumar/ai-issue-agent/internal/repository"
	"github.com/yashashkumar/ai-issue-agent/internal/service"
	"github.com/yashashkumar/ai-issue-agent/internal/spawner"
)

type WebhookHandler struct {
	projectRepo *repository.ProjectRepo
	spawner     *spawner.Spawner
}

func NewWebhookHandler(projectRepo *repository.ProjectRepo, spawner *spawner.Spawner) *WebhookHandler {
	return &WebhookHandler{
		projectRepo: projectRepo,
		spawner:     spawner,
	}
}

func (h *WebhookHandler) HandleGitHubIssues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)

	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	p, err := h.projectRepo.GetByGitHubRepo(ctx, owner, repo)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "project not found", reqID)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to read body", reqID)
		return
	}

	if p.GitHubWebhookSecret != nil && *p.GitHubWebhookSecret != "" {
		sig := r.Header.Get("X-Hub-Signature-256")
		if err := service.VerifySignature(*p.GitHubWebhookSecret, sig, bodyBytes); err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid signature", reqID)
			return
		}
	}

	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "issues" {
		writeJSON(w, http.StatusOK, map[string]bool{"ignored": true})
		return
	}

	payload, err := service.ParseWebhookPayload(bodyBytes)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json payload", reqID)
		return
	}

	validAction := false
	for _, act := range []string{"opened", "edited", "reopened", "assigned"} {
		if payload.Action == act {
			validAction = true
			break
		}
	}

	if !validAction {
		writeJSON(w, http.StatusOK, map[string]bool{"ignored": true})
		return
	}

	emailAllowed := false
	senderEmail := payload.Issue.User.Email
	if senderEmail == "" {
		// Mock logic or check sender login directly if email missing
		senderEmail = payload.Sender.Login + "@github.com"
	}

	for _, e := range p.AllowedEmails {
		if e == senderEmail || e == payload.Sender.Login {
			emailAllowed = true
			break
		}
	}

	// Also allow if allowed_emails is empty (permissive configuration) or logic config
	if len(p.AllowedEmails) > 0 && !emailAllowed {
		writeJSONError(w, http.StatusForbidden, "sender not in allowed list", reqID)
		return // Forbidden
	}

	agentID := uuid.NewString()

	var systemPrompt string
	if payload.Action == "opened" {
		systemPrompt = fmt.Sprintf("Title: %s\n\nDescription: %s", payload.Issue.Title, payload.Issue.Body)
	} else {
		systemPrompt = fmt.Sprintf("Title: %s\n\nDescription: %s\n\n---\nContext: This issue was %s. Continue working on the existing branch.",
			payload.Issue.Title, payload.Issue.Body, payload.Action)
	}

	issueNumber := payload.Issue.Number

	req := spawner.SpawnRequest{
		ID:                agentID,
		Source:            models.AgentSourceGitHubIssue,
		ProjectID:         &p.ID,
		SystemPrompt:      systemPrompt,
		WorkDir:           fmt.Sprintf("%s/worktrees/issue-%d", p.RootFolder, issueNumber),
		Name:              fmt.Sprintf("Issue #%d Agent", issueNumber),
		Description:       fmt.Sprintf("Agent spawned for GitHub issue #%d (%s)", issueNumber, payload.Action),
		GitHubIssueNumber: &issueNumber,
		GitHubOwner:       &owner,
		GitHubRepo:        &repo,
		CreateBranch:      payload.Action == "opened",
		CreatePR:          true,
	}

	h.spawner.Spawn(req)

	writeJSON(w, http.StatusAccepted, map[string]string{
		"agent_id": agentID,
		"status":   "pending",
	})
}

// Additional HandleGitHubIssueComments logic
func (h *WebhookHandler) HandleGitHubIssueComments(w http.ResponseWriter, r *http.Request) {
	// similar parsing Logic ...
	ctx := r.Context()
	reqID := GetRequestID(ctx)

	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	p, err := h.projectRepo.GetByGitHubRepo(ctx, owner, repo)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "project not found", reqID)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to read body", reqID)
		return
	}

	if p.GitHubWebhookSecret != nil && *p.GitHubWebhookSecret != "" {
		sig := r.Header.Get("X-Hub-Signature-256")
		if err := service.VerifySignature(*p.GitHubWebhookSecret, sig, bodyBytes); err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid signature", reqID)
			return
		}
	}

	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "issue_comment" {
		writeJSON(w, http.StatusOK, map[string]bool{"ignored": true})
		return
	}

	payload, err := service.ParseWebhookPayload(bodyBytes)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json payload", reqID)
		return
	}

	if payload.Action != "created" && payload.Action != "edited" {
		writeJSON(w, http.StatusOK, map[string]bool{"ignored": true})
		return
	}

	agentID := uuid.NewString()
	issueNumber := payload.Issue.Number

	systemPrompt := fmt.Sprintf("Issue Title: %s\n\nIssue Description: %s\n\n---\nNew Comment by @%s:\n%s\n\n---\nInstructions: Address the feedback/request in the comment above. You are working on branch agent/issue-%d.",
		payload.Issue.Title, payload.Issue.Body, payload.Sender.Login, payload.Comment.Body, issueNumber)

	req := spawner.SpawnRequest{
		ID:                agentID,
		Source:            models.AgentSourceGitHubComment,
		ProjectID:         &p.ID,
		SystemPrompt:      systemPrompt,
		WorkDir:           fmt.Sprintf("%s/worktrees/issue-%d", p.RootFolder, issueNumber),
		Name:              fmt.Sprintf("Comment Agent Issue #%d", issueNumber),
		Description:       "Feedback handling",
		GitHubIssueNumber: &issueNumber,
		GitHubOwner:       &owner,
		GitHubRepo:        &repo,
		CreateBranch:      false,
		CreatePR:          true,
	}

	h.spawner.Spawn(req)

	writeJSON(w, http.StatusAccepted, map[string]string{
		"agent_id": agentID,
		"status":   "pending",
	})
}
