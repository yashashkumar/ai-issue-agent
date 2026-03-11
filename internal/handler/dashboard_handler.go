package handler

import (
	"log/slog"
	"net/http"

	"github.com/yashashkumar/ai-issue-agent/internal/service"
	"github.com/yashashkumar/ai-issue-agent/internal/web/templates/pages"
)

type DashboardHandler struct {
	agentSvc   *service.AgentService
	projectSvc *service.ProjectService
	logger     *slog.Logger
}

func NewDashboardHandler(agentSvc *service.AgentService, projectSvc *service.ProjectService, logger *slog.Logger) *DashboardHandler {
	return &DashboardHandler{
		agentSvc:   agentSvc,
		projectSvc: projectSvc,
		logger:     logger,
	}
}

func (h *DashboardHandler) Home(w http.ResponseWriter, r *http.Request) {
	// Add context passing if needed
	err := pages.Home().Render(r.Context(), w)
	if err != nil {
		h.logger.Error("failed to render home page", "err", err)
	}
}

// Stubs for the rest of the dashboard
func (h *DashboardHandler) ProjectsList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Projects list not fully implemented. Run via API."))
}

func (h *DashboardHandler) ProjectFormNew(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Project form not fully implemented. Create via API."))
}
