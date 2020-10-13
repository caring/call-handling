package handlers

import (
	"context"

	"github.com/caring/call-handling/internal/db"
	"github.com/caring/call-handling/pb"
	"github.com/caring/go-packages/pkg/errors"
	"google.golang.org/grpc/codes"
)

type callMethods interface {
	Create(context.Context, *db.Call) error
}

func CreateCall(ctx context.Context, in *pb.CreateCallRequest, store callMethods) (resp *pb.CallhandlingResponse, err error) {
	var (
		call *db.Call
	)
	call = db.NewCall(in)
	err = store.Create(ctx, call)
	if err != nil {
		err = errors.WithGrpcStatus(err, codes.Internal)
	}
	resp = call.ToProto()
	return
}
