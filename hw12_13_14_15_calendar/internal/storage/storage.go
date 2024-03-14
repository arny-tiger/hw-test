package storage

import (
	"context"

	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage/entity"
)

const (
	SQLStorageType    = "sql"
	MemoryStorageType = "memory"
)

type EventNotFoundErr struct {
	Msg string
}

func (e *EventNotFoundErr) Error() string {
	return e.Msg
}

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	CreateEvent(evt entity.Event) (int, error)
	UpdateEvent(evt entity.Event) (int, error)
	DeleteEvent(id int) error
	GetEvents(limit int, offset int) ([]entity.Event, error)
	GetEvent(id int) (*entity.Event, error)
}
