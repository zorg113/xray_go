package memorystorage

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[int64]*storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[int64]*storage.Event),
	}
}

func (s *Storage) CreateEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.events[e.ID] != nil {
		var msg strings.Builder
		fmt.Fprintf(&msg, "evant %d already exist", e.ID)
		logrus.Info(msg)
		return errors.New("event already exist")
	}
	s.events[e.ID] = &e
	return nil
}

func (s *Storage) UpdateEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] == nil {
		var msg strings.Builder
		fmt.Fprintf(&msg, "event %d not found", e.ID)
		logrus.Info(msg)
		return errors.New("event not found")
	}

	s.events[e.ID] = &e
	return nil
}

func (s *Storage) DeleteEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] == nil {
		var msg strings.Builder
		fmt.Fprintf(&msg, "event %d not found", e.ID)
		logrus.Info(msg)
		return errors.New("event not found")
	}

	delete(s.events, e.ID)
	return nil
}

func (s *Storage) GetEvents(startData time.Time, endData time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var events []storage.Event
	for _, event := range s.events {
		if event.StartDate.Before(endData) && event.EndDate.After(startData) {
			events = append(events, *event)
		}
	}
	return events, nil
}
