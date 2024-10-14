package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amedoeyes/hadath/internal/api/response"
	"github.com/amedoeyes/hadath/internal/service"
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
	err := h.service.Create(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *BookingHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.service.ListByCurrentUser(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) ListByEvent(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.service.ListByEvent(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.service.Delete(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
