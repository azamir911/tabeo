package services

import (
	"space_trouble_booking/internal/models"
	"space_trouble_booking/internal/repository"
)

type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) CreateBooking(booking *models.Booking) error {
	return s.repo.CreateBooking(booking)
}

func (s *BookingService) GetAllBookings() ([]models.Booking, error) {
	return s.repo.GetAllBookings()
}
