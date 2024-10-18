package main

import (
	"database/sql"
	"log"
	"net/http"

	"space_trouble_booking/internal/handlers"
	"space_trouble_booking/internal/repository"
	"space_trouble_booking/internal/services"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to the database
	db, err := sql.Open("postgres", "postgres://postgres:postgres@db:5432/space_trouble?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repository, service, and handler
	repo := repository.NewBookingRepository(db)
	service := services.NewBookingService(repo)
	handler := handlers.NewBookingHandler(service)

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/bookings", handler.CreateBooking).Methods("POST")
	router.HandleFunc("/bookings", handler.GetAllBookings).Methods("GET")

	// Start the server
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
