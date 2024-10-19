package main

import (
	"database/sql"
	"fmt"
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

	// Initialize services and handlers
	repo := repository.NewBookingRepository(db)
	spaceXService := services.NewSpaceXService()
	bookingService := services.NewBookingService(repo, spaceXService)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	// Set up the router
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is up and running!")
	}).Methods("GET")

	router.HandleFunc("/bookings", bookingHandler.CreateBooking).Methods("POST")
	router.HandleFunc("/bookings", bookingHandler.GetAllBookings).Methods("GET")

	// Start the server
	log.Println("Startingggg server on :8080...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
