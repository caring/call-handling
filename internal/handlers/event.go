package handlers

import (
	"context"

	"google.golang.org/grpc/codes"

	"github.com/caring/call-handling/internal/db"
	"github.com/caring/call-handling/pb"
	"github.com/caring/go-packages/pkg/errors"
)

const (
	DIAL       = "dialing"
	RING       = "ringing"
	CONNECT    = "connected"
	DISCONNECT = "disconnected"
	JOIN       = "party joined"
	EXIT       = "party exited"
	DISPO      = "dispositioned"
)

type eventMethods interface {
	Create(context.Context, *db.Event) error
}

func Dialed(ctx context.Context, in *pb.EventRequest, store eventMethods) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, store, DIAL)
}

func Ringed(ctx context.Context, in *pb.EventRequest, store eventMethods) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, store, RING)
}

func Connected(ctx context.Context, in *pb.EventRequest, store eventMethods) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, store, CONNECT)
}

func Disconnected(ctx context.Context, in *pb.EventRequest, store eventMethods) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, store, DISCONNECT)
}

func Joined(ctx context.Context, in *pb.EventRequest, store eventMethods) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, store, JOIN)
}

func Exited(ctx context.Context, in *pb.EventRequest, store eventMethods) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, store, EXIT)
}

func Dispositioned(ctx context.Context, in *pb.EventRequest, store eventMethods) (*pb.CallhandlingResponse, error) {
	return createEvent(ctx, in, store, DISPO)
}

func createEvent(ctx context.Context, in *pb.EventRequest, store eventMethods, eventType string) (resp *pb.CallhandlingResponse, err error) {
	var (
		event *db.Event
	)
	event = db.NewEvent(in, eventType)
	err = store.Create(ctx, event)
	if err != nil {
		err = errors.WithGrpcStatus(err, codes.Internal)
	}
	resp = event.ToProto()
	return
}
