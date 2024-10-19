package services

import (
	"errors"
	"space_trouble_booking/internal/models"
	"space_trouble_booking/internal/repository"
	"time"
)

type BookingService struct {
	repo          *repository.BookingRepository
	spaceXService *SpaceXService
}

func NewBookingService(repo *repository.BookingRepository, spaceXService *SpaceXService) *BookingService {
	return &BookingService{
		repo:          repo,
		spaceXService: spaceXService,
	}
}

func (s *BookingService) CreateBooking(booking *models.Booking) error {
	// Check for launchpad availability with SpaceX data
	if err := s.validateLaunchpadAvailability(booking.LaunchpadID, booking.LaunchDate); err != nil {
		return err
	}

	// Proceed to create the booking if no conflicts are found
	return s.repo.CreateBooking(booking)
}

func (s *BookingService) GetAllBookings() ([]models.Booking, error) {
	return s.repo.GetAllBookings()
}

// validateLaunchpadAvailability checks if the launchpad is available for the given date
func (s *BookingService) validateLaunchpadAvailability(launchpadID string, launchDate time.Time) error {
	launches, err := s.spaceXService.GetUpcomingLaunches()
	if err != nil {
		return err
	}

	for _, launch := range launches {
		// If the launchpad ID and date match, it means there's a conflict
		if launch.Launchpad == launchpadID && launch.DateUTC.Truncate(24*time.Hour).Equal(launchDate.Truncate(24*time.Hour)) {
			return errors.New("the selected launchpad is not available on this date due to a SpaceX launch")
		}
	}

	return nil
}
