package main

import (
	"context"
	"fmt"
	"time"

	"github.com/caring/call-handling/internal/db"

	// _ "github.com/caring/call-handling/internal/handlers"
	"github.com/caring/call-handling/pb"
	_ "github.com/caring/go-packages/pkg/errors"
	_ "google.golang.org/grpc/codes"
)

const (
	DIAL       = "dialing"
	RING       = "ringing"
	CONNECT    = "connected"
	DISCONNECT = "disconnected"
	JOIN       = "party joined"
	EXIT       = "party exited"
	DISPO      = "dispositioned"
	ENQUEUE    = "enqueued"
	VOICEMAIL  = "voicemail created"
)

type service struct {
}

func (s *service) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	l.Info(fmt.Sprintf("Received: %v", in.Data))
	resp := "Data: " + in.Data

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	status := "up"
	if err := store.Ping(ctx); err != nil {
		status = "down"
	}
	return &pb.PingResponse{Data: resp + "; Database: " + status}, nil

}

func (s *service) CreateCall(ctx context.Context, in *pb.CreateCallRequest) (*pb.CallhandlingResponse, error) {
	l.Info(fmt.Sprintf("Received: %v", in.Call))

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := store.Calls.Create(ctx, db.NewCall(in))

	return &pb.CallhandlingResponse{}, err
}

func (s *service) DialEvent(ctx context.Context, in *pb.EventRequest) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, DIAL)
}

func (s *service) RingEvent(ctx context.Context, in *pb.EventRequest) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, RING)
}

func (s *service) ConnectEvent(ctx context.Context, in *pb.EventRequest) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, CONNECT)
}

func (s *service) DisconnectEvent(ctx context.Context, in *pb.EventRequest) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, DISCONNECT)
}

func (s *service) JoinEvent(ctx context.Context, in *pb.EventRequest) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, JOIN)
}

func (s *service) ExitEvent(ctx context.Context, in *pb.EventRequest) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, EXIT)
}

func (s *service) DispositionEvent(ctx context.Context, in *pb.EventRequest) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, DISPO)
}

func (s *service) EnqueuedEvent(ctx context.Context, in *pb.EventRequest) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, ENQUEUE)
}

func (s *service) VoicemailEvent(ctx context.Context, in *pb.EventRequest) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, VOICEMAIL)
}

func createEvent(ctx context.Context, in *pb.EventRequest, eventType string) (*pb.CallhandlingResponse, error) {
	l.Info(fmt.Sprintf("Received: %v", in.Event))
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	err := store.Events.Create(ctx, db.NewEvent(in, eventType))
	return &pb.CallhandlingResponse{}, err
}
