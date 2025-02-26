package storage

import "time"

type Event struct {
	ID          string
	Title       string
	StartData   time.Time
	EndData     time.Time
	Description string
	OwnerID     string
	RemindIn    string
}
