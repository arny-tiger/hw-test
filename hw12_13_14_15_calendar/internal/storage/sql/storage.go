package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres
)

type Storage struct {
	db *sqlx.DB
}

type updateRequestFields struct {
	name  string
	value interface{}
}

func New(config config.Config) (*Storage, error) {
	dsn := fmt.Sprintf(
		"sslmode=disable host=%s port=%d dbname=%s user=%s password=%s",
		config.DB.Host, config.DB.Port, config.DB.Database, config.DB.Username, config.DB.Password,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	err := s.db.PingContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateEvent(evt entity.Event) (int, error) {
	query := "INSERT INTO events (title, date, duration, description, owner_id)" +
		" VALUES (:title, :date, :duration, :description, :owner_id) RETURNING id"
	rows, err := s.db.NamedQuery(query, evt)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var lastInsertID int
	if rows.Next() {
		err := rows.Scan(&lastInsertID)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("CreateEvent failed")
	}
	return lastInsertID, nil
}

func (s *Storage) UpdateEvent(evt entity.Event) (int, error) {
	if evt.ID == 0 {
		return 0, fmt.Errorf("UpdateEvent failed. No ID")
	}
	query := "UPDATE events SET"
	params := map[string]interface{}{}
	fields := []updateRequestFields{
		{"title", evt.Title},
		{"date", evt.Date},
		{"duration", evt.Duration},
		{"description", evt.Description},
		{"owner_id", evt.OwnerID},
	}
	for _, field := range fields {
		if field.value != "" && field.value != 0 {
			query += fmt.Sprintf(" %s = :%s", field.name, field.name)
			params[field.name] = field.value
		}
	}
	if len(params) == 0 {
		return 0, fmt.Errorf("UpdateEvent failed. No fields to change")
	}

	query += " WHERE id = :id"
	params["id"] = evt.ID
	_, err := s.db.NamedExec(query, params)
	if err != nil {
		return 0, err
	}
	return evt.ID, nil
}

func (s *Storage) DeleteEvent(id int) error {
	query := "DELETE FROM events WHERE id = $1"
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetEvents(limit int, offset int) ([]entity.Event, error) {
	query := "SELECT * FROM events LIMIT $1 OFFSET $2"
	var events []entity.Event
	err := s.db.Select(&events, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Storage) GetEvent(id int) (*entity.Event, error) {
	query := "SELECT * FROM events WHERE id = $1"
	var event entity.Event
	err := s.db.Get(&event, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &storage.EventNotFoundErr{Msg: err.Error()}
		}
		return nil, err
	}
	return &event, nil
}
