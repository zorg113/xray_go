package storage

import "time"

type Event struct {
	ID          int64
	Title       string
	StartDate   time.Time
	EndDate     time.Time
	Description string
	OwnerID     string
	RemindIn    string
}
