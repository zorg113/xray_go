package memorystorage

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	mu     sync.RWMutex //nolint:unused
	events map[string]*storage.Event
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.events[e.ID] != nil {
		var msg strings.Builder
		fmt.Fprintf(&msg, "evant %s already exist", e.ID)
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
		fmt.Fprintf(&msg, "event %s not found", e.ID)
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
		fmt.Fprintf(&msg, "event %s not found", e.ID)
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
		if event.StartData.Second() >= startData.Second() && event.EndData.Second() <= endData.Second() {
			events = append(events, *event)
		}
	}
	return events, nil
}
