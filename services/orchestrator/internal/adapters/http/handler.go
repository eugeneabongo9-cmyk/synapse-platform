package http

import (
	"encoding/json"
	"net/http"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/application"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Handler struct {
	service *application.OrchestratorService
	logger  zerolog.Logger
}

func NewHandler(service *application.OrchestratorService, logger zerolog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/agents", h.createAgent)
	mux.HandleFunc("POST /api/v1/agents/{id}/start", h.startAgent)
	mux.HandleFunc("GET /health", h.health)

	return mux
}

func (h *Handler) createAgent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	agent, err := h.service.CreateAgent(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":     agent.ID.String(),
		"name":   agent.Name,
		"status": agent.Status,
	})
}

func (h *Handler) startAgent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid agent id", http.StatusBadRequest)
		return
	}

	if err := h.service.StartAgent(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
