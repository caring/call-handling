package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/caring/call-handling/pb"
	"github.com/caring/go-packages/pkg/errors"
)

type eventService struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

type Event struct {
	CallID     int64
	Type       string
	IdentityID int64
	Timestamp  int64
	Meta       string
}

type protoEvent interface {
	GetEvent() *pb.Event
}

func NewEvent(proto protoEvent, eventType string) *Event {
	e := proto.GetEvent()
	return &Event{
		CallID:     e.GetCallId(),
		Type:       eventType,
		IdentityID: e.GetIdentityId(),
		Timestamp:  e.GetTimestamp(),
		Meta:       e.GetMeta(),
	}
}

func (m *Event) ToProto() *pb.EventResponse {
	return &pb.EventResponse{
		CallId:     m.CallID,
		IdentityId: m.IdentityID,
		Timestamp:  m.Timestamp,
		Meta:       m.Meta,
		Type:       m.Type,
	}
}

func (svc *eventService) Get(ctx context.Context, ID int64) (*Event, error) {
	return svc.get(ctx, false, ID)
}

func (svc *eventService) GetTx(ctx context.Context, ID int64) (*Event, error) {
	return svc.get(ctx, true, ID)
}

func (svc *eventService) get(ctx context.Context, useTx bool, ID int64) (*Event, error) {
	errMsg := func() string { return "Error executing get event - " + fmt.Sprint(ID) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return nil, err
		}

		stmt = tx.Stmt(svc.stmts["get-event"])
	} else {
		stmt = svc.stmts["get-event"]
	}

	p := Event{}

	err = stmt.QueryRowContext(ctx, ID).
		Scan(&p.CallID, &p.Type, &p.IdentityID, &p.Timestamp, &p.Meta)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(ErrNotFound, errMsg())
		}

		return nil, errors.Wrap(err, errMsg())
	}

	return &p, nil
}

func (svc *eventService) Create(ctx context.Context, input *Event) error {
	return svc.create(ctx, false, input)
}

func (svc *eventService) CreateTx(ctx context.Context, input *Event) error {
	return svc.create(ctx, true, input)
}

func (svc *eventService) create(ctx context.Context, useTx bool, input *Event) error {
	errMsg := func() string { return "Error executing create event - " + fmt.Sprint(input) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["create-event"])
	} else {
		stmt = svc.stmts["create-event"]
	}

	result, err := stmt.ExecContext(ctx, input.CallID, input.Type, input.IdentityID, input.Timestamp, input.Meta)
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	if rowCount == 0 {
		return errors.Wrap(ErrNoRowsAffected, errMsg())
	}

	return nil
}
