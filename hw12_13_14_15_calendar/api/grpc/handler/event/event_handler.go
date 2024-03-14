package event

//go:generate protoc --proto_path=../../proto --go_out=../../gen --go-grpc_out=../../gen ../../proto/EventService.proto

import (
	"context"
	"fmt"

	pbEvents "github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/api/grpc/gen/event"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage/entity"
)

type Handler struct {
	pbEvents.UnimplementedEventServiceServer
	Storage storage.Storage
	Logger  logger.Logger
}

func (s *Handler) CreateEvent(_ context.Context, req *pbEvents.CreateEventRequest) (*pbEvents.CreateEventResponse, error) {
	eventEntity := convertPBMessageToEntity(req.GetEvent())
	eventID, err := s.Storage.CreateEvent(eventEntity)
	if err != nil {
		err = fmt.Errorf("CreateEvent failed: %w", err)
		s.Logger.Error(err.Error())
		return nil, err
	}
	event, err := s.Storage.GetEvent(eventID)
	if err != nil {
		err = fmt.Errorf("CreateEvent failed. Unable to get event with ID %d: %w", eventID, err)
		s.Logger.Error(err.Error())
		return nil, err
	}
	pbEvent := convertEntityToPBMessage(*event)
	return &pbEvents.CreateEventResponse{Event: &pbEvent}, nil
}

func (s *Handler) GetEvents(_ context.Context, req *pbEvents.GetEventsRequest) (*pbEvents.GetEventsResponse, error) {
	limit, offset := req.GetLimit(), req.GetOffset()
	if limit == 0 {
		limit = 20
	}
	dbEvents, err := s.Storage.GetEvents(int(limit), int(offset))
	if err != nil {
		return nil, err
	}
	events := make([]*pbEvents.Event, 0, len(dbEvents))
	for _, event := range dbEvents {
		pbEvent := convertEntityToPBMessage(event)
		events = append(events, &pbEvent)
	}
	return &pbEvents.GetEventsResponse{Events: events}, nil
}

func (s *Handler) UpdateEvent(_ context.Context, req *pbEvents.UpdateEventRequest) (*pbEvents.UpdateEventResponse, error) {
	eventEntity := convertPBMessageToEntity(req.GetEvent())
	eventID, err := s.Storage.UpdateEvent(eventEntity)
	if err != nil {
		err = fmt.Errorf("UpdateEvent failed: %w", err)
		s.Logger.Error(err.Error())
		return nil, err
	}
	event, err := s.Storage.GetEvent(eventID)
	if err != nil {
		err = fmt.Errorf("UpdateEvent failed. Unable to get event with ID %d: %w", eventID, err)
		s.Logger.Error(err.Error())
		return nil, err
	}
	pbEvent := convertEntityToPBMessage(*event)
	return &pbEvents.UpdateEventResponse{Event: &pbEvent}, nil
}

func (s *Handler) DeleteEvent(_ context.Context, req *pbEvents.DeleteEventRequest) (*pbEvents.DeleteEventResponse, error) {
	eventID := req.GetId()
	_, err := s.Storage.GetEvent(int(eventID))
	if err != nil {
		err = fmt.Errorf("DeleteEvent failed. Unable to get event with ID %d: %w", eventID, err)
		s.Logger.Error(err.Error())
		return nil, err
	}
	err = s.Storage.DeleteEvent(int(eventID))
	if err != nil {
		err = fmt.Errorf("DeleteEvent failed. Unable to delete event with ID %d: %w", eventID, err)
		s.Logger.Error(err.Error())
		return nil, err
	}
	return &pbEvents.DeleteEventResponse{ResponseCode: 1}, nil
}

func convertPBMessageToEntity(eventMsg *pbEvents.Event) entity.Event {
	return entity.Event{
		ID:          int(eventMsg.GetID()),
		Title:       eventMsg.GetTitle(),
		Date:        eventMsg.GetDate(),
		Duration:    eventMsg.GetDuration(),
		Description: eventMsg.GetDescription(),
		OwnerID:     int(eventMsg.GetOwnerID()),
	}
}

func convertEntityToPBMessage(event entity.Event) pbEvents.Event {
	return pbEvents.Event{
		ID:          int32(event.ID),
		Title:       event.Title,
		Date:        event.Date,
		Duration:    event.Duration,
		Description: event.Description,
		OwnerID:     int32(event.OwnerID),
	}
}
