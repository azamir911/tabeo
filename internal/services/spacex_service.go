package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SpaceXService struct {
	client *http.Client
}

func NewSpaceXService() *SpaceXService {
	return &SpaceXService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Launch represents a simplified structure of a SpaceX launch.
type Launch struct {
	Launchpad string    `json:"launchpad"`
	DateUTC   time.Time `json:"date_utc"`
}

// GetUpcomingLaunches fetches upcoming SpaceX launches from the API.
func (s *SpaceXService) GetUpcomingLaunches() ([]Launch, error) {
	url := "https://api.spacexdata.com/v4/launches/upcoming"
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from SpaceX API, status code: %d", resp.StatusCode)
	}

	var launches []Launch
	if err := json.NewDecoder(resp.Body).Decode(&launches); err != nil {
		return nil, err
	}

	return launches, nil
}
