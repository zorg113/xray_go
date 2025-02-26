package app

import (
	"time"

	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Storage interface {
	CreateEvent(e storage.Event) error
	UpdateEvent(e storage.Event) error
	DeleteEvent(e storage.Event) error
	GetEvents(startData, endData time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(e storage.Event) error {
	if err := a.storage.CreateEvent(e); err != nil {
		a.logger.Error("cann't create event")
		return err
	}
	return nil
}

func (a *App) UpdateEvent(e storage.Event) error {
	if err := a.storage.UpdateEvent(e); err != nil {
		a.logger.Error("can't update event")
		return err
	}
	return nil
}

func (a *App) DeleteEvent(e storage.Event) error {
	if err := a.storage.DeleteEvent(e); err != nil {
		a.logger.Error("can't delete event")
		return err
	}
	return nil
}

func (a *App) GetEvents(startDate, endDate time.Time) ([]storage.Event, error) {
	var events []storage.Event
	var err error
	if events, err = a.storage.GetEvents(startDate, endDate); err != nil {
		a.logger.Error("can't delete event")
		return nil, err
	}
	return events, nil
}
