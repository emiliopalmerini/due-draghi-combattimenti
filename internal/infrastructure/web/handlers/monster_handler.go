package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"

	monsterApp "github.com/emiliopalmerini/due-draghi-combattimenti/internal/application/monster"
	"github.com/emiliopalmerini/due-draghi-combattimenti/internal/infrastructure/web/templates"
)

// MonsterHandler handles HTTP requests for monster browsing.
type MonsterHandler struct {
	service *monsterApp.Service
	logger  *slog.Logger
}

// NewMonsterHandler creates a new monster HTTP handler.
func NewMonsterHandler(service *monsterApp.Service, logger *slog.Logger) *MonsterHandler {
	return &MonsterHandler{
		service: service,
		logger:  logger,
	}
}

// SearchHandler handles monster search requests via HTMX.
// GET /api/monsters?max_xp=N&q=search
func (h *MonsterHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())

	maxXP := 1_000_000
	if v := r.URL.Query().Get("max_xp"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			h.logger.Error("Invalid max_xp", "request_id", requestID, "error", err)
			http.Error(w, "Invalid max_xp parameter", http.StatusBadRequest)
			return
		}
		maxXP = parsed
	}

	query := r.URL.Query().Get("q")
	monsters := h.service.SearchMonsters(query, maxXP)

	w.Header().Set("Content-Type", "text/html")
	if err := templates.MonsterList(monsters, maxXP).Render(r.Context(), w); err != nil {
		h.logger.Error("Failed to render monster list", "request_id", requestID, "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
