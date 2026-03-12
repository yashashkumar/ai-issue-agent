package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yashashkumar/ai-issue-agent/internal/models"
	"github.com/yashashkumar/ai-issue-agent/internal/service"
)

type ProjectHandler struct {
	projectSvc *service.ProjectService
}

func NewProjectHandler(projectSvc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectSvc: projectSvc}
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)

	// In a real app we'd parse ?search=...
	search := r.URL.Query().Get("search")

	var projects []models.Project
	var err error

	if search != "" {
		// Use a search method on the service... wait, I'll add search logic directly or via service
		// The service currently only has ListProjects. Let's just return ListProjects for now,
		// or if you want, fetch all and filter in memory since this is a simple implementation.
		all, err2 := h.projectSvc.ListProjects(ctx)
		if err2 != nil {
			err = err2
		} else {
			for _, p := range all {
				if p.Name == search || p.Description == search { // very naive search
					projects = append(projects, p)
				}
			}
		}
	} else {
		projects, err = h.projectSvc.ListProjects(ctx)
	}

	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "failed to list projects", reqID)
		return
	}

	if projects == nil {
		projects = []models.Project{}
	}

	WriteJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)

	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid request body", reqID)
		return
	}

	if p.Name == "" || p.RootFolder == "" {
		WriteJSONError(w, http.StatusBadRequest, "name and root_folder are required", reqID)
		return
	}

	created, err := h.projectSvc.CreateProject(ctx, &p)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "failed to create project", reqID)
		return
	}

	WriteJSON(w, http.StatusCreated, created)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)
	id := r.PathValue("id")

	p, err := h.projectSvc.GetProject(ctx, id)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, "project not found", reqID)
		return
	}

	WriteJSON(w, http.StatusOK, p)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)
	id := r.PathValue("id")

	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid request body", reqID)
		return
	}
	p.ID = id

	if err := h.projectSvc.UpdateProject(ctx, &p); err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "failed to update project", reqID)
		return
	}

	updated, _ := h.projectSvc.GetProject(ctx, id)
	WriteJSON(w, http.StatusOK, updated)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetRequestID(ctx)
	id := r.PathValue("id")

	if err := h.projectSvc.DeleteProject(ctx, id); err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "failed to delete project", reqID)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
