package repository

import (
	"database/sql"
	"space_trouble_booking/internal/models"
	"time"
)

type H2BookingRepository struct {
	db *sql.DB
}

func NewH2BookingRepository(db *sql.DB) *H2BookingRepository {
	return &H2BookingRepository{db: db}
}

func (r *H2BookingRepository) Migrate() error {
	createTable := `
        CREATE TABLE IF NOT EXISTS bookings (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            first_name TEXT,
            last_name TEXT,
            gender TEXT,
            birthday TEXT,
            launchpad_id TEXT,
            destination TEXT,
            launch_date TEXT
        );
    `
	_, err := r.db.Exec(createTable)
	return err
}

func (r *H2BookingRepository) CreateBooking(booking *models.Booking) (*models.Booking, error) {
	query := `
        INSERT INTO bookings (first_name, last_name, gender, birthday, launchpad_id, destination, launch_date)
        VALUES (?, ?, ?, ?, ?, ?, ?);
    `
	result, err := r.db.Exec(query, booking.FirstName, booking.LastName, booking.Gender, booking.Birthday, booking.LaunchpadID, booking.Destination, booking.LaunchDate)
	if err != nil {
		return nil, err
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Update the booking ID with the newly generated one
	booking.ID = int(id)

	return booking, nil
}

func (r *H2BookingRepository) GetAllBookings() ([]*models.Booking, error) {
	rows, err := r.db.Query("SELECT id, first_name, last_name, gender, birthday, launchpad_id, destination, launch_date FROM bookings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var booking models.Booking
		var birthdayStr, launchDateStr string

		err := rows.Scan(&booking.ID, &booking.FirstName, &booking.LastName, &booking.Gender, &birthdayStr, &booking.LaunchpadID, &booking.Destination, &launchDateStr)
		if err != nil {
			return nil, err
		}

		// Parse the birthday and launch date strings into time.Time using the correct format
		layout := "2006-01-02 15:04:05-07:00"
		birthday, err := time.Parse(layout, birthdayStr)
		if err != nil {
			return nil, err
		}
		launchDate, err := time.Parse(layout, launchDateStr)
		if err != nil {
			return nil, err
		}

		booking.Birthday = birthday
		booking.LaunchDate = launchDate

		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

func (r *H2BookingRepository) FindBookingByLaunchpadAndDate(launchpadID string, launchDate time.Time) (*models.Booking, error) {
	var booking models.Booking
	var birthdayStr, launchDateStr string

	// Query to find a booking by launchpad ID and date
	query := `
        SELECT id, first_name, last_name, gender, birthday, launchpad_id, destination, launch_date
        FROM bookings
        WHERE launchpad_id = ? AND DATE(launch_date) = DATE(?);
    `
	err := r.db.QueryRow(query, launchpadID, launchDate.Format("2006-01-02 15:04:05-07:00")).Scan(
		&booking.ID, &booking.FirstName, &booking.LastName, &booking.Gender,
		&birthdayStr, &booking.LaunchpadID, &booking.Destination, &launchDateStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil // No conflict found
	}
	if err != nil {
		return nil, err
	}

	// Parse the birthday and launch date from the string format
	parsedBirthday, err := time.Parse("2006-01-02 15:04:05-07:00", birthdayStr)
	if err != nil {
		return nil, err
	}
	parsedLaunchDate, err := time.Parse("2006-01-02 15:04:05-07:00", launchDateStr)
	if err != nil {
		return nil, err
	}
	booking.Birthday = parsedBirthday
	booking.LaunchDate = parsedLaunchDate

	return &booking, nil
}

func (r *H2BookingRepository) DeleteBooking(id int) error {
	query := "DELETE FROM bookings WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
