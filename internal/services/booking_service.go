package services

import (
	"errors"
	"space_trouble_booking/internal/models"
	"space_trouble_booking/internal/repository"
	"time"
)

type BookingService struct {
	repo          repository.BookingRepository
	spaceXService *SpaceXService
}

func NewBookingService(repo repository.BookingRepository, spaceXService *SpaceXService) *BookingService {
	return &BookingService{
		repo:          repo,
		spaceXService: spaceXService,
	}
}

func (s *BookingService) CreateBooking(booking *models.Booking) (*models.Booking, error) {
	// Check for launchpad availability with SpaceX data
	if err := s.validateLaunchpadAvailability(booking.LaunchpadID, booking.LaunchDate); err != nil {
		return nil, err
	}

	// Proceed to create the booking if no conflicts are found
	createdBooking, err := s.repo.CreateBooking(booking)
	if err != nil {
		return nil, err
	}

	return createdBooking, nil
}

func (s *BookingService) GetAllBookings() ([]*models.Booking, error) {
	return s.repo.GetAllBookings()
}

func (s *BookingService) DeleteBooking(id int) error {
	return s.repo.DeleteBooking(id)
}

// validateLaunchpadAvailability checks if the launchpad is available for the given date
func (s *BookingService) validateLaunchpadAvailability(launchpadID string, launchDate time.Time) error {
	// Check SpaceX launches for conflicts
	launches, err := s.spaceXService.GetUpcomingLaunches()
	if err != nil {
		return err
	}

	for _, launch := range launches {
		launchTime := launch.DateUTC

		// Compare only the date parts by truncating to the beginning of the day
		if launch.Launchpad == launchpadID && launchTime.Truncate(24*time.Hour).Equal(launchDate.Truncate(24*time.Hour)) {
			return errors.New("the selected launchpad is not available on this date due to a SpaceX launch")
		}
	}

	// Check internal bookings for conflicts
	existingBooking, err := s.repo.FindBookingByLaunchpadAndDate(launchpadID, launchDate)
	if err != nil {
		return err
	}
	if existingBooking != nil {
		return errors.New("the selected launchpad is not available on this date due to an existing booking")
	}

	return nil
}
