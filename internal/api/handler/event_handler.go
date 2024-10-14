package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amedoeyes/hadath/internal/api/request"
	"github.com/amedoeyes/hadath/internal/api/response"
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
	var req request.EventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	err = h.service.Create(r.Context(), &req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.List(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	event := h.service.Get(r.Context())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req request.EventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	err = h.service.Update(r.Context(), &req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.service.Delete(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
