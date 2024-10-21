package repository

import (
	"database/sql"
	"space_trouble_booking/internal/models"
	"time"
)

type PostgresBookingRepository struct {
	db *sql.DB
}

func NewPostgresBookingRepository(db *sql.DB) *PostgresBookingRepository {
	return &PostgresBookingRepository{db: db}
}

func (r *PostgresBookingRepository) Migrate() error {
	createTable := `
        CREATE TABLE IF NOT EXISTS bookings (
            id SERIAL PRIMARY KEY,
            first_name TEXT,
            last_name TEXT,
            gender TEXT,
            birthday TIMESTAMPTZ,
            launchpad_id TEXT,
            destination TEXT,
            launch_date TIMESTAMPTZ
        );
    `
	_, err := r.db.Exec(createTable)
	return err
}

func (r *PostgresBookingRepository) CreateBooking(booking *models.Booking) (*models.Booking, error) {
	query := `
        INSERT INTO bookings (first_name, last_name, gender, birthday, launchpad_id, destination, launch_date)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;
    `

	// Convert birthday and launch date to RFC3339 string format
	birthdayStr := booking.Birthday.Format(time.RFC3339)
	launchDateStr := booking.LaunchDate.Format(time.RFC3339)

	var id int
	err := r.db.QueryRow(query, booking.FirstName, booking.LastName, booking.Gender, birthdayStr, booking.LaunchpadID, booking.Destination, launchDateStr).Scan(&id)
	if err != nil {
		return nil, err
	}

	booking.ID = id
	return booking, nil
}

func (r *PostgresBookingRepository) GetAllBookings() ([]*models.Booking, error) {
	rows, err := r.db.Query("SELECT id, first_name, last_name, gender, birthday, launchpad_id, destination, launch_date FROM bookings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var booking models.Booking
		var birthdayStr, launchDateStr string

		// Scan fields, storing date fields as strings
		err := rows.Scan(&booking.ID, &booking.FirstName, &booking.LastName, &booking.Gender,
			&birthdayStr, &booking.LaunchpadID, &booking.Destination, &launchDateStr)
		if err != nil {
			return nil, err
		}

		// Parse birthday and launch_date from strings to time.Time
		birthday, err := time.Parse(time.RFC3339, birthdayStr)
		if err != nil {
			return nil, err
		}
		launchDate, err := time.Parse(time.RFC3339, launchDateStr)
		if err != nil {
			return nil, err
		}

		booking.Birthday = birthday
		booking.LaunchDate = launchDate

		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

func (r *PostgresBookingRepository) FindBookingByLaunchpadAndDate(launchpadID string, launchDate time.Time) (*models.Booking, error) {
	var booking models.Booking
	var launchDateStr string

	// Query to find a booking by launchpad ID and date
	query := `
        SELECT id, first_name, last_name, gender, birthday, launchpad_id, destination, launch_date
        FROM bookings
        WHERE launchpad_id = $1 AND DATE(launch_date) = DATE($2);
    `
	err := r.db.QueryRow(query, launchpadID, launchDate.Format("2006-01-02 15:04:05-07:00")).Scan(
		&booking.ID, &booking.FirstName, &booking.LastName, &booking.Gender,
		&booking.Birthday, &booking.LaunchpadID, &booking.Destination, &launchDateStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil // No conflict found
	}
	if err != nil {
		return nil, err
	}

	// Parse the launch date from the string format
	parsedLaunchDate, err := time.Parse("2006-01-02 15:04:05-07:00", launchDateStr)
	if err != nil {
		return nil, err
	}
	booking.LaunchDate = parsedLaunchDate

	return &booking, nil
}

func (r *PostgresBookingRepository) DeleteBooking(id int) error {
	query := "DELETE FROM bookings WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
