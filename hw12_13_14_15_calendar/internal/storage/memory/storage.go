package memorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/entity"
)

type Storage struct {
	events    map[int]entity.Event
	idCounter int
	mu        sync.RWMutex
}

type EventNotFoundErr struct {
	Msg string
}

func (e *EventNotFoundErr) Error() string {
	return e.Msg
}

func New() *Storage {
	return &Storage{
		events: make(map[int]entity.Event),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateEvent(evt entity.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.idCounter++
	evt.ID = s.idCounter
	s.events[evt.ID] = evt
	return evt.ID, nil
}

func (s *Storage) UpdateEvent(evt entity.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[evt.ID]; !ok {
		return 0, &EventNotFoundErr{fmt.Sprintf("event with id %d not found", evt.ID)}
	}
	s.events[evt.ID] = evt
	return evt.ID, nil
}

func (s *Storage) DeleteEvent(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return &EventNotFoundErr{fmt.Sprintf("event with id %d not found", id)}
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) GetEvents() ([]entity.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]entity.Event, 0, len(s.events))
	for _, evt := range s.events {
		events = append(events, evt)
	}
	return events, nil
}
