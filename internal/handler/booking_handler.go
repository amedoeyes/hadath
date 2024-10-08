package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amedoeyes/hadath/internal/model"
	"github.com/amedoeyes/hadath/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{
		service: service,
	}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request struct {
		EventID uuid.UUID `json:"event_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(*model.User)

	err = h.service.Create(r.Context(), user.ID, request.EventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *BookingHandler) GetAllByUserID(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*model.User)

	bookings, err := h.service.GetAllByUserID(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) GetAllByEventID(w http.ResponseWriter, r *http.Request) {
	event := r.Context().Value("event").(*model.Event)

	bookings, err := h.service.GetAllByEventID(r.Context(), event.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.DeleteByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
