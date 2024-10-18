package repository

import (
	"database/sql"

	"space_trouble_booking/internal/models"

	_ "github.com/lib/pq" // Postgres driver
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) CreateBooking(booking *models.Booking) error {
	query := `
        INSERT INTO bookings (first_name, last_name, gender, birthday, launchpad_id, destination, launch_date)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
    `
	return r.db.QueryRow(query,
		booking.FirstName,
		booking.LastName,
		booking.Gender,
		booking.Birthday,
		booking.LaunchpadID,
		booking.Destination,
		booking.LaunchDate,
	).Scan(&booking.ID)
}

func (r *BookingRepository) GetAllBookings() ([]models.Booking, error) {
	var bookings []models.Booking
	query := `SELECT * FROM bookings`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.ID, &b.FirstName, &b.LastName, &b.Gender, &b.Birthday, &b.LaunchpadID, &b.Destination, &b.LaunchDate); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}
