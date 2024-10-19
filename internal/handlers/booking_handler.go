package handlers

import (
	"encoding/json"
	"net/http"
	"space_trouble_booking/internal/models"
	"space_trouble_booking/internal/services"
	"time"
)

type BookingHandler struct {
	service *services.BookingService
}

func NewBookingHandler(service *services.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Parse launch date from string to time.Time
	launchDate, err := time.Parse("2006-01-02", booking.LaunchDate.Format("2006-01-02"))
	if err != nil {
		http.Error(w, "Invalid launch date", http.StatusBadRequest)
		return
	}
	booking.LaunchDate = launchDate

	if err := h.service.CreateBooking(&booking); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.service.GetAllBookings()
	if err != nil {
		http.Error(w, "Failed to retrieve bookings", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookings)
}
