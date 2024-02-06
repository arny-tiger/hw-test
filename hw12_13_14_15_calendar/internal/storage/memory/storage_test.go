package memorystorage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	s := New()
	ctx := context.Background()
	err := s.Connect(ctx)
	require.NoError(t, err)
}

func TestClose(t *testing.T) {
	s := New()
	err := s.Close()
	require.NoError(t, err)
}

func TestCreateEventConcurrent(t *testing.T) {
	s := New()
	eventsAmnt := 30
	wg := sync.WaitGroup{}
	wg.Add(eventsAmnt)

	e := entity.Event{
		Title:       "Event Title",
		Date:        time.Now().String(),
		Duration:    "",
		Description: "Event Description",
		OwnerID:     1,
	}
	for i := 0; i < eventsAmnt; i++ {
		go func() {
			defer wg.Done()
			_, err := s.CreateEvent(e)
			if err != nil {
				require.NoError(t, err)
			}
		}()
	}
	wg.Wait()
	lastInsertedID := s.events[len(s.events)].ID
	require.Equal(t, eventsAmnt, lastInsertedID)
}

func TestUpdateEventConcurrent(t *testing.T) {
	existedEvent := entity.Event{
		ID:          1,
		Title:       "Event Title",
		Date:        time.Now().String(),
		Duration:    "",
		Description: "Event Description",
		OwnerID:     1,
	}
	s := New()
	s.events[existedEvent.ID] = existedEvent

	updatedEntity := entity.Event{
		ID:          1,
		Title:       "New Event Title",
		Date:        time.Now().AddDate(0, 0, 1).String(),
		Duration:    "",
		Description: "New Event Description",
		OwnerID:     1,
	}

	updatesAmount := 15
	wg := sync.WaitGroup{}
	wg.Add(updatesAmount)

	for i := 0; i < updatesAmount; i++ {
		go func() {
			defer wg.Done()
			_, err := s.UpdateEvent(updatedEntity)
			require.NoError(t, err)
		}()
	}
	wg.Wait()
	require.Equal(t, updatedEntity, s.events[existedEvent.ID])
}

func TestDeleteEventConcurrent(t *testing.T) {
	s := New()
	e := entity.Event{
		ID:          1,
		Title:       "Event Title",
		Date:        time.Now().String(),
		Duration:    "",
		Description: "Event Description",
		OwnerID:     1,
	}
	s.events[e.ID] = e

	deletesAmount := 15
	wg := sync.WaitGroup{}
	wg.Add(deletesAmount)

	for i := 0; i < deletesAmount; i++ {
		go func() {
			defer wg.Done()
			err := s.DeleteEvent(e.ID)
			if err != nil {
				_, ok := s.events[e.ID]
				if !ok {
					require.Error(t, &EventNotFoundErr{})
				}
			}
		}()
	}
	wg.Wait()
	require.Empty(t, s.events)
}

func TestGetEventsConcurrent(t *testing.T) {
	s := New()
	event1 := entity.Event{
		ID:          1,
		Title:       "Event Title",
		Date:        time.Now().String(),
		Duration:    "",
		Description: "Event Description",
		OwnerID:     1,
	}
	event2 := entity.Event{
		ID:          2,
		Title:       "Event Title 2",
		Date:        time.Now().String(),
		Duration:    "",
		Description: "Event Description 2",
		OwnerID:     2,
	}
	s.events[event1.ID] = event1
	s.events[event2.ID] = event2

	requestsAmount := 15
	wg := sync.WaitGroup{}
	wg.Add(requestsAmount)

	for i := 0; i < requestsAmount; i++ {
		go func() {
			defer wg.Done()
			events, err := s.GetEvents()
			require.NoError(t, err)
			require.ElementsMatch(t, events, []entity.Event{event1, event2})
		}()
	}
	wg.Wait()
}
