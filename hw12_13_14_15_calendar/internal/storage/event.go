package storage

import (
	"time"
)

type Event struct {
	ID          uint64
	Title       string
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Description string
	OwnerID     string
	RemindIn    string
}
