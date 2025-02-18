package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
)

type Row struct {
	ID          int64
	Title       string
	StartDate   time.Time
	EndDate     time.Time
	Decsription string
	OwnerID     int64
	RemindIn    int64
}

type Storage struct {
	db *sql.DB
}

func New(ctx context.Context, user, passwrd, host, name string, port uint64) (*Storage, error) {
	config := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable", user, passwrd, host, port, name)
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
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) NewEvent(e storage.Event) error {
	_, err := s.db.Exec(
		`INSERT INTO event (title, start_date, end_date, description, owner_id, remind_in)`+
			` VALUES $1, $2, $3, $4, $5, $6 RETURNING id`,
		e.Title,
		e.StartData,
		e.EndData,
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
		e.StartData,
		e.EndData,
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

func (s *Storage) GetEvents(startData time.Time, endData time.Time) ([]Row, error) {
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

	var row Row
	var rows []Row
	for events.Next() {
		if err := events.Scan(&row.ID, &row.Title, &row.StartDate, &row.EndDate, &row.Decsription, &row.OwnerID, &row.RemindIn); err != nil {
			log.Fatal(err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}
