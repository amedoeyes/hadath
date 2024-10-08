package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/amedoeyes/hadath/internal/model"
	"github.com/amedoeyes/hadath/internal/service"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Address     string    `json:"address"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(*model.User)

	err := h.service.Create(
		r.Context(),
		user.ID,
		request.Name,
		request.Description,
		request.Address,
		request.StartTime,
		request.EndTime,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *EventHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	var response struct {
		ID          uint32    `json:"id"`
		UserID      uint32    `json:"user_id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Address     string    `json:"address"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	event := r.Context().Value("event").(*model.Event)

	response.ID = event.ID
	response.UserID = event.UserID
	response.Name = event.Name
	response.Description = event.Description
	response.Address = event.Address
	response.StartTime = event.StartTime
	response.EndTime = event.EndTime
	response.CreatedAt = event.CreatedAt
	response.UpdatedAt = event.UpdatedAt

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *EventHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        *string    `json:"name,omitempty"`
		Description *string    `json:"description,omitempty"`
		Address     *string    `json:"address,omitempty"`
		StartTime   *time.Time `json:"start_time,omitempty"`
		EndTime     *time.Time `json:"end_time,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := r.Context().Value("event").(*model.Event)

	if request.Name != nil {
		event.Name = *request.Name
	}
	if request.Description != nil {
		event.Description = *request.Description
	}
	if request.Address != nil {
		event.Address = *request.Address
	}
	if request.StartTime != nil {
		event.StartTime = *request.StartTime
	}
	if request.EndTime != nil {
		event.EndTime = *request.EndTime
	}

	err := h.service.UpdateByID(
		r.Context(),
		event.ID,
		event.Name,
		event.Description,
		event.Address,
		event.StartTime,
		event.EndTime,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *EventHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	event := r.Context().Value("event").(*model.Event)

	err := h.service.DeleteByID(r.Context(), event.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
