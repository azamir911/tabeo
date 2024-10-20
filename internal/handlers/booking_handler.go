package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"space_trouble_booking/internal/models"
	"space_trouble_booking/internal/services"
	"strconv"
)

type BookingHandler struct {
	service *services.BookingService
}

func NewBookingHandler(service *services.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking

	// Decode the JSON request body into the booking struct
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		log.Printf("Got a request error: '%+v'", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service to create the booking
	createdBooking, err := h.service.CreateBooking(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// Return the created booking as a response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdBooking)
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

func (h *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	// Extract the booking ID from the URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing booking ID", http.StatusBadRequest)
		return
	}

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// Call the service to delete the booking
	err = h.service.DeleteBooking(id)
	if err != nil {
		http.Error(w, "Failed to delete booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content indicates successful deletion
}
