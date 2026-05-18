package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

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
		h.log.Error("Failed to deserialize body", "error", err)
		http.Error(w, "Failed to deserialize body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.CreateDepartment(department); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(department)
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	splitURL := strings.Split(r.URL.Path, "/")

	departmentID, err := strconv.Atoi(splitURL[len(splitURL)-3])
	if err != nil {
		h.log.Error("Invalid department ID", "error", err)
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	var req models.RequestEmployee

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to deserialize body", "error", err)
		http.Error(w, "Failed to deserialize body", http.StatusBadRequest)
		return
	}

	var hiredAt *time.Time
	if req.HiredAt != "" {
		t, err := time.Parse("02-01-2006", req.HiredAt)
		if err != nil {
			http.Error(w, "Invalid department ID, use DD-MM-YYYY format", http.StatusBadRequest)
			return
		}
		hiredAt = &t
	}

	employee := &models.Employee{
		DepartmentId: departmentID,
		FullName:     req.FullName,
		Position:     req.Position,
		HiredAt:      hiredAt,
	}

	if err = h.usecase.CreateEmployee(employee); err != nil {
		h.log.Error("Failed to create employee", "error", err)
		http.Error(w, "Failed to create employee", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
}

func (h *Handler) ExistingDepartments(w http.ResponseWriter, r *http.Request) {

	idParam := path.Base(r.URL.Path)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.log.Error("Invalid department ID", "error", err)
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}
	h.log.Debug("from URL", "id", id)

	switch r.Method {
	case http.MethodGet:

		var req models.RequestTree

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("Failed to deserialize body", "error", err)
			http.Error(w, "Failed to deserialize body", http.StatusBadRequest)
			return
		}

		req.Id = id

		department, err := h.usecase.GetTree(&req)
		if err != nil {
			h.log.Error("Failed to get department tree", "error", err)
			http.Error(w, "Failed to get department tree", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(department)

	case http.MethodPatch:

		var department *models.Department
		if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
			h.log.Error("Failed to deserialize body", "error", err)
			http.Error(w, "Failed to deserialize body", http.StatusBadRequest)
			return
		}

		department.Id = id

		if err = h.usecase.UpdateParent(department); err != nil {
			if strings.Contains(err.Error(), "cannot make department subtree of its subtree") {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(department)

	case http.MethodDelete:

		var req models.RequestDelete
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("Failed to deserialize body", "error", err)
			http.Error(w, "Failed to deserialize body", http.StatusBadRequest)
			return
		}

		req.Id = id

		if err := h.usecase.DeleteDepartment(&req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
