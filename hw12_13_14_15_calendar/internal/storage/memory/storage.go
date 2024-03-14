package memorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage/entity"
)

type Storage struct {
	events    map[int]entity.Event
	idCounter int
	mu        sync.RWMutex
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
		return 0, &storage.EventNotFoundErr{Msg: fmt.Sprintf("event with id %d not found", evt.ID)}
	}
	s.events[evt.ID] = evt
	return evt.ID, nil
}

func (s *Storage) DeleteEvent(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return &storage.EventNotFoundErr{Msg: fmt.Sprintf("event with id %d not found", id)}
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) GetEvents(limit int, offset int) ([]entity.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if offset > len(s.events) {
		return []entity.Event{}, nil
	}

	end := limit + offset
	if end > len(s.events) {
		end = len(s.events)
	}

	events := make([]entity.Event, 0, len(s.events))
	for _, evt := range s.events {
		events = append(events, evt)
	}
	return events[offset:end], nil
}

func (s *Storage) GetEvent(id int) (*entity.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.events[id]; !ok {
		return nil, &storage.EventNotFoundErr{Msg: fmt.Sprintf("event with id %d not found", id)}
	}
	event := s.events[id]
	return &event, nil
}
