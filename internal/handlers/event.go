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

func Dialed(ctx context.Context, in *pb.EventRequest, store eventMethods) (resp *pb.CallhandlingResponse, err error) {
	var (
		event *db.Event
	)
	event = db.NewEvent(in, DIAL)
	err = store.Create(ctx, event)
	if err != nil {
		err = errors.WithGrpcStatus(err, codes.Internal)
	}
	resp = event.ToProto()
	return
}

func Ringed(ctx context.Context, in *pb.EventRequest, store eventMethods) (resp *pb.CallhandlingResponse, err error) {
	var (
		event *db.Event
	)
	event = db.NewEvent(in, RING)
	err = store.Create(ctx, event)
	if err != nil {
		err = errors.WithGrpcStatus(err, codes.Internal)
	}
	resp = event.ToProto()
	return
}

func Connected(ctx context.Context, in *pb.EventRequest, store eventMethods) (resp *pb.CallhandlingResponse, err error) {
	var (
		event *db.Event
	)
	event = db.NewEvent(in, CONNECT)
	err = store.Create(ctx, event)
	if err != nil {
		err = errors.WithGrpcStatus(err, codes.Internal)
	}
	resp = event.ToProto()
	return
}

func Disconnected(ctx context.Context, in *pb.EventRequest, store eventMethods) (resp *pb.CallhandlingResponse, err error) {
	var (
		event *db.Event
	)
	event = db.NewEvent(in, DISCONNECT)
	err = store.Create(ctx, event)
	if err != nil {
		err = errors.WithGrpcStatus(err, codes.Internal)
	}
	resp = event.ToProto()
	return
}

func Joined(ctx context.Context, in *pb.EventRequest, store eventMethods) (resp *pb.CallhandlingResponse, err error) {
	var (
		event *db.Event
	)
	event = db.NewEvent(in, JOIN)
	err = store.Create(ctx, event)
	if err != nil {
		err = errors.WithGrpcStatus(err, codes.Internal)
	}
	resp = event.ToProto()
	return
}

func Exited(ctx context.Context, in *pb.EventRequest, store eventMethods) (resp *pb.CallhandlingResponse, err error) {
	var (
		event *db.Event
	)
	event = db.NewEvent(in, EXIT)
	err = store.Create(ctx, event)
	if err != nil {
		err = errors.WithGrpcStatus(err, codes.Internal)
	}
	resp = event.ToProto()
	return
}

func Dispositioned(ctx context.Context, in *pb.EventRequest, store eventMethods) (resp *pb.CallhandlingResponse, err error) {
	var (
		event *db.Event
	)
	event = db.NewEvent(in, DISPO)
	err = store.Create(ctx, event)
	if err != nil {
		err = errors.WithGrpcStatus(err, codes.Internal)
	}
	resp = event.ToProto()
	return
}
