package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/caring/go-packages/pkg/errors"

	"github.com/caring/call-handling/pb"
)

// callService provides an API for interacting with the calls table
type callService struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

// Call is a struct representation of a row in the calls table
type Call struct {
	ID             int64
	SID            int64
	ConversationID int64
	ANI            string
	DNIS           string
	Status         string
}

// protoCall is an interface that most proto call objects will satisfy
type protoCall interface {
	GetCall() *pb.Call
}

// NewCall is a convenience helper cast a proto call to it's DB layer struct
func NewCall(proto protoCall) *Call {
	c := proto.GetCall()
	return &Call{
		ID:             c.GetCallId(),
		SID:            c.GetSid(),
		ConversationID: c.GetConversationId(),
		ANI:            c.GetANI(),
		DNIS:           c.GetDNIS(),
		Status:         c.GetStatus(),
	}
}

// ToProto casts a db call into a proto response object
func (m *Call) ToProto() *pb.CallResponse {
	return &pb.CallResponse{
		CallId:         m.ID,
		Sid:            m.SID,
		ConversationId: m.ConversationID,
		ANI:            m.ANI,
		DNIS:           m.DNIS,
		Status:         m.Status,
	}
}

// Get fetches a single call from the db
func (svc *callService) Get(ctx context.Context, ID int64) (*Call, error) {
	return svc.get(ctx, false, ID)
}

// GetTx fetches a single call from the db inside of a tx from ctx
func (svc *callService) GetTx(ctx context.Context, ID int64) (*Call, error) {
	return svc.get(ctx, true, ID)
}

// get fetches a single call from the db
func (svc *callService) get(ctx context.Context, useTx bool, ID int64) (*Call, error) {
	errMsg := func() string { return "Error executing get call - " + fmt.Sprint(ID) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return nil, err
		}

		stmt = tx.Stmt(svc.stmts["get-call"])
	} else {
		stmt = svc.stmts["get-call"]
	}

	p := Call{}

	err = stmt.QueryRowContext(ctx, ID).
		Scan(&p.ID, &p.SID, &p.ConversationID, &p.ANI, &p.DNIS, &p.Status)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(ErrNotFound, errMsg())
		}

		return nil, errors.Wrap(err, errMsg())
	}

	return &p, nil
}

// Create a new call
func (svc *callService) Create(ctx context.Context, input *Call) error {
	return svc.create(ctx, false, input)
}

// CreateTx creates a new call withing a tx from ctx
func (svc *callService) CreateTx(ctx context.Context, input *Call) error {
	return svc.create(ctx, true, input)
}

// create a new call. if useTx = true then it will attempt to create the callhandling within a transaction
// from context.
func (svc *callService) create(ctx context.Context, useTx bool, input *Call) error {
	errMsg := func() string { return "Error executing create call - " + fmt.Sprint(input) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["create-call"])
	} else {
		stmt = svc.stmts["create-call"]
	}

	result, err := stmt.ExecContext(ctx, input.ID, input.SID, input.ConversationID, input.ANI, input.DNIS, input.Status)
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

// Update updates a single call row in the DB
func (svc *callService) Update(ctx context.Context, input *Call) error {
	return svc.update(ctx, false, input)
}

// UpdateTx updates a single call row in the DB within a tx from ctx
func (svc *callService) UpdateTx(ctx context.Context, input *Call) error {
	return svc.update(ctx, true, input)
}

// update a call. if useTx = true then it will attempt to update the callhandling within a transaction
// from context.
func (svc *callService) update(ctx context.Context, useTx bool, input *Call) error {
	errMsg := func() string { return "Error executing update callhandling - " + fmt.Sprint(input) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["update-call"])
	} else {
		stmt = svc.stmts["update-call"]
	}

	result, err := stmt.ExecContext(ctx, input.ID, input.SID, input.ConversationID, input.ANI, input.DNIS, input.Status)
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

// Delete sets deleted_at for a single calls row
func (svc *callService) Delete(ctx context.Context, ID int64) error {
	return svc.delete(ctx, false, ID)
}

// DeleteTx sets deleted_at for a single calls row within a tx from ctx
func (svc *callService) DeleteTx(ctx context.Context, ID int64) error {
	return svc.delete(ctx, true, ID)
}

// delete a callhandling by setting deleted at. if useTx = true then it will attempt to delete the callhandling within a transaction
// from context.
func (svc *callService) delete(ctx context.Context, useTx bool, ID int64) error {
	errMsg := func() string { return "Error executing delete call - " + fmt.Sprint(ID) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["delete-call"])
	} else {
		stmt = svc.stmts["delete-call"]
	}

	result, err := stmt.ExecContext(ctx, ID)
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
