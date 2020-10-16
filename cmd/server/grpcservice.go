package main

import (
	"context"
	"fmt"
	"time"

	"github.com/caring/call-handling/internal/handlers"

	"github.com/caring/call-handling/pb"
	_ "github.com/caring/go-packages/pkg/errors"
	_ "google.golang.org/grpc/codes"
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

func (s *service) CreateCall(ctx context.Context, in *pb.CallRequest) (*pb.CallResponse, error) {
	return handlers.CreateCall(ctx, in, store.Calls)
}

func (s *service) Dialed(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	return handlers.Dialed(ctx, in, store.Events)
}

func (s *service) Ringed(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	return handlers.Ringed(ctx, in, store.Events)
}

func (s *service) Connected(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	return handlers.Connected(ctx, in, store.Events)
}

func (s *service) Disconnected(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	return handlers.Disconnected(ctx, in, store.Events)
}

func (s *service) Joined(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	return handlers.Joined(ctx, in, store.Events)
}

func (s *service) Exited(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	return handlers.Exited(ctx, in, store.Events)
}

func (s *service) Dispositioned(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	return handlers.Dispositioned(ctx, in, store.Events)
}

func (s *service) Enqueued(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	return handlers.Enqueued(ctx, in, store.Events)

}
