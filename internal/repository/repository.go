package repository

import (
	"space_trouble_booking/internal/models"
	"time"
)

type BookingRepository interface {
	Migrate() error
	CreateBooking(booking *models.Booking) (*models.Booking, error)
	GetAllBookings() ([]*models.Booking, error)
	FindBookingByLaunchpadAndDate(launchpadID string, launchDate time.Time) (*models.Booking, error)
}
