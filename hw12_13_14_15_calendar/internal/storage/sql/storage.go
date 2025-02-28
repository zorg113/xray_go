package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/config"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(ctx context.Context, conf config.DBConf) (*Storage, error) {
	config := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName)
	db, err := sql.Open("postgres", config)
	if err != nil {
		return nil, fmt.Errorf("cannot connect db: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot ping context: %w", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	_, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	err = s.db.PingContext(ctx)
	if err != nil {
		return errors.Wrap(err, "connot ping context")
	}
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(e storage.Event) error {
	_, err := s.db.Exec(
		`INSERT INTO event (title, start_date, end_date, description, owner_id, remind_in)`+
			` VALUES $1, $2, $3, $4, $5, $6 RETURNING id`,
		e.Title,
		e.StartDate,
		e.EndDate,
		e.Description,
		e.OwnerID,
		e.RemindIn,
	)
	if err != nil {
		return errors.Wrap(err, "cannot create event")
	}
	return nil
}

func (s *Storage) UpdateEvent(e storage.Event) error {
	_, err := s.db.Exec(
		`UPDATE event SET (`+
			`title, start_date, end_date, description, owner_id, remind_in`+
			`) = (`+
			`$1, $2, $3, $4, $5, $6`+
			`) WHERE id = $7`,
		e.Title,
		e.StartDate,
		e.EndDate,
		e.Description,
		e.OwnerID,
		e.RemindIn,
		e.ID,
	)
	if err != nil {
		return errors.Wrap(err, "cannot update event")
	}
	return nil
}

func (s *Storage) DeleteEvent(e storage.Event) error {
	_, err := s.db.Exec(
		`DELETE FROM event WHERE id=$1`,
		e.ID,
	)
	if err != nil {
		return errors.Wrap(err, "cannot delete event")
	}
	return nil
}

func (s *Storage) GetEvents(startData time.Time, endData time.Time) ([]storage.Event, error) {
	events, err := s.db.Query(
		`SELECT id, 
       			title, 
       			start_date, 
    		    end_date, 
    		    description, 
    		    owner_id, 
    		    remind_in
			FROM event
			WHERE start_date >=$1 AND start_date <=$2`,
		startData,
		endData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "cannot execute query")
	}
	defer events.Close()

	var row storage.Event
	var rows []storage.Event
	for events.Next() {
		if err := events.Scan(
			&row.ID,
			&row.Title,
			&row.StartDate,
			&row.EndDate,
			&row.Description,
			&row.OwnerID,
			&row.RemindIn); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}
