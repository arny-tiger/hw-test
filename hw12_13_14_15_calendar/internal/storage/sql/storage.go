package sqlstorage

import (
	"context"
	"fmt"

	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres
)

type Storage struct {
	db *sqlx.DB
}

func New(config config.Config) (*Storage, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
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
		" VALUES (:title, :date, :duration, :description, :owner_id)"
	res, err := s.db.NamedExec(query, evt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *Storage) UpdateEvent(evt entity.Event) (int, error) {
	query := "UPDATE events SET title = :title, date = :date, duration = :duration," +
		" description = :description, owner_id = :owner_id WHERE id = :id"
	res, err := s.db.NamedExec(query, evt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *Storage) DeleteEvent(id int) error {
	query := "DELETE FROM events WHERE id = :id"
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetEvents() ([]entity.Event, error) {
	query := "SELECT * FROM events"
	var events []entity.Event
	err := s.db.Select(&events, query)
	if err != nil {
		return nil, err
	}
	return events, nil
}
