package models

import "time"

type Booking struct {
	ID          int       `db:"id" json:"id"`
	FirstName   string    `db:"first_name" json:"first_name"`
	LastName    string    `db:"last_name" json:"last_name"`
	Gender      string    `db:"gender" json:"gender"`
	Birthday    time.Time `db:"birthday" json:"birthday"`
	LaunchpadID string    `db:"launchpad_id" json:"launchpad_id"`
	Destination string    `db:"destination" json:"destination"`
	LaunchDate  time.Time `db:"launch_date" json:"launch_date"`
}
