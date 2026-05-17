package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/VysMax/organizational-structure/models"
	"github.com/VysMax/organizational-structure/usecase"
)

type Handler struct {
	usecase *usecase.Usecase
	log     *slog.Logger
}

func New(uc *usecase.Usecase, log *slog.Logger) *Handler {
	return &Handler{
		usecase: uc,
		log:     log,
	}
}

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var department *models.Department
	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		http.Error(w, "Failed to deserialize body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.CreateDepartment(department); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(department)
}
