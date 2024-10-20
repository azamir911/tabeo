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
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

func main() {
	// Connect to the H2 (SQLite) database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use the H2 implementation of the BookingRepository
	repo := repository.NewH2BookingRepository(db)
	// Migrate the database schema
	if err := repo.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

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
	router.HandleFunc("/bookings/{id}", bookingHandler.DeleteBooking).Methods("DELETE")

	// Start the server
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
