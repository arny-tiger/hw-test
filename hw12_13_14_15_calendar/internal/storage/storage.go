package storage

import (
	"context"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/entity"
)

const (
	SQLStorageType    = "sql"
	MemoryStorageType = "memory"
)

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	CreateEvent(evt entity.Event) (int, error)
	UpdateEvent(evt entity.Event) (int, error)
	DeleteEvent(id int) error
	GetEvents() ([]entity.Event, error)
}
