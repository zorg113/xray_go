package grpc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/storage"
	cnd_pb "github.com/zorg113/xray_go/hw12_13_14_15_calendar/pkg/calendar"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//nolint:all
func (a *Server) CreateEvent(_ context.Context, req *cnd_pb.CreateEventRequest) (*cnd_pb.CreateEventResponse, error) {
	event := convertToStorageEvent(req.Event)
	if err := a.app.CreateEvent(event); err != nil {
		return nil, fmt.Errorf("can't create event: %w", err)
	}
	return &cnd_pb.CreateEventResponse{}, nil
}

func (a *Server) UpdateEvent(_ context.Context, req *cnd_pb.UpdateEventRequest) (*cnd_pb.UpdateEventResponse, error) {
	event := convertToStorageEvent(req.Event)
	if err := a.app.UpdateEvent(event); err != nil {
		return nil, fmt.Errorf("can't update event: %w", err)
	}
	return &cnd_pb.UpdateEventResponse{}, nil
}

func (a *Server) GetEvents(_ context.Context, req *cnd_pb.GetEventsRequest) (*cnd_pb.GetEventsResponse, error) {
	events, err := a.app.GetEvents(req.StartDate.AsTime(), req.EndDate.AsTime())
	if err != nil {
		return nil, fmt.Errorf("can't get events: %w", err)
	}
	var pbEvents []*cnd_pb.Event //nolint: prealloc
	for _, val := range events {
		var pbEvent *cnd_pb.Event
		pbEvent, err = convertToProtoBufEvent(val)
		if err != nil {
			return nil, fmt.Errorf("can't convert to event: %w", err)
		}
		pbEvents = append(pbEvents, pbEvent)
	}
	return &cnd_pb.GetEventsResponse{Event: pbEvents}, nil
}

func (a *Server) DeleteEvent(_ context.Context, req *cnd_pb.DeleteEventRequest) (*cnd_pb.DeleteEventsResponse, error) {
	event := convertToStorageEvent(req.Event)
	if err := a.app.DeleteEvent(event); err != nil {
		return nil, fmt.Errorf("can't delete event: %w", err)
	}
	return &cnd_pb.DeleteEventsResponse{}, nil
}

func convertToStorageEvent(e *cnd_pb.Event) storage.Event {
	return storage.Event{
		ID:          e.Id,
		Title:       e.Title,
		StartDate:   e.StartDate.AsTime(),
		EndDate:     e.EndDate.AsTime(),
		Description: e.Description,
		OwnerID:     fmt.Sprintf("%d", e.OwnerId),
		RemindIn:    e.RemindIn,
	}
}

func convertToProtoBufEvent(e storage.Event) (*cnd_pb.Event, error) {
	id := e.ID
	ownerID, err := strconv.Atoi(e.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("can't convert to int: %w", err)
	}
	return &cnd_pb.Event{
		Id:          id,
		Title:       e.Title,
		StartDate:   &timestamppb.Timestamp{Seconds: int64(e.StartDate.Second())},
		EndDate:     &timestamppb.Timestamp{Seconds: int64(e.EndDate.Second())},
		Description: e.Description,
		OwnerId:     uint64(ownerID), //nolint:gosec
		RemindIn:    e.RemindIn,
	}, nil
}
