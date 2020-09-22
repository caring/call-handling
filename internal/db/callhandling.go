package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/caring/go-packages/pkg/errors"
	"github.com/google/uuid"

	"github.com/caring/call-handling/pb"
)



// callhandlingService provides an API for interacting with the callhandlings table
type callhandlingService struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

// Callhandling is a struct representation of a row in the callhandlings table
type Callhandling struct {
	ID  	uuid.UUID
	Name  string
}

// protoCallhandling is an interface that most proto callhandling objects will satisfy
type protoCallhandling interface {
	GetName() string
}

// NewCallhandling is a convenience helper cast a proto callhandling to it's DB layer struct
func NewCallhandling(ID string, proto protoCallhandling) (*Callhandling, error) {
	mID, err := ParseUUID(ID)
	if err != nil {
		return nil, err
	}

	return &Callhandling{
		ID:  	mID,
		Name: proto.GetName(),
	}, nil
}

// ToProto casts a db callhandling into a proto response object
func (m *Callhandling) ToProto() *pb.CallhandlingResponse {
	return &pb.CallhandlingResponse{
		Id:  				m.ID.String(),
		Name:       m.Name,
	}
}

// Get fetches a single callhandling from the db
func (svc *callhandlingService) Get(ctx context.Context, ID uuid.UUID) (*Callhandling, error) {
	return svc.get(ctx, false, ID)
}

// GetTx fetches a single callhandling from the db inside of a tx from ctx
func (svc *callhandlingService) GetTx(ctx context.Context, ID uuid.UUID) (*Callhandling, error) {
	return svc.get(ctx, true, ID)
}

// get fetches a single callhandling from the db
func (svc *callhandlingService) get(ctx context.Context, useTx bool, ID uuid.UUID) (*Callhandling, error) {
	errMsg := func() string { return "Error executing get callhandling - " + fmt.Sprint(ID) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return nil, err
		}

		stmt = tx.Stmt(svc.stmts["get-callhandling"])
	} else {
		stmt = svc.stmts["get-callhandling"]
	}

	p := Callhandling{}

	err = stmt.QueryRowContext(ctx, ID).
		Scan(&m.CallhandlingID, &m.Name)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(ErrNotFound, errMsg())
		}

		return nil, errors.Wrap(err, errMsg())
	}

	return &p, nil
}

// Create a new callhandling
func (svc *callhandlingService) Create(ctx context.Context, input *Callhandling) error {
	return svc.create(ctx, false, input)
}

// CreateTx creates a new callhandling withing a tx from ctx
func (svc *callhandlingService) CreateTx(ctx context.Context, input *Callhandling) error {
	return svc.create(ctx, true, input)
}

// create a new callhandling. if useTx = true then it will attempt to create the callhandling within a transaction
// from context.
func (svc *callhandlingService) create(ctx context.Context, useTx bool, input *Callhandling) error {
	errMsg := func() string { return "Error executing create callhandling - " + fmt.Sprint(input) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["create-callhandling"])
	} else {
		stmt = svc.stmts["create-callhandling"]
	}

	result, err := stmt.ExecContext(ctx, input.CallhandlingID, input.Name)
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	if rowCount == 0 {
		return errors.Wrap(ErrNotCreated, errMsg())
	}

	return nil
}

// Update updates a single callhandling row in the DB
func (svc *callhandlingService) Update(ctx context.Context, input *Callhandling) error {
	return svc.update(ctx, false, input)
}

// UpdateTx updates a single callhandling row in the DB within a tx from ctx
func (svc *callhandlingService) UpdateTx(ctx context.Context, input *Callhandling) error {
	return svc.update(ctx, true, input)
}

// update a callhandling. if useTx = true then it will attempt to update the callhandling within a transaction
// from context.
func (svc *callhandlingService) update(ctx context.Context, useTx bool, input *Callhandling) error {
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

		stmt = tx.Stmt(svc.stmts["update-callhandling"])
	} else {
		stmt = svc.stmts["update-callhandling"]
	}

	result, err := stmt.ExecContext(ctx, input.Name, input.CallhandlingID)
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

// Delete sets deleted_at for a single callhandlings row
func (svc *callhandlingService) Delete(ctx context.Context, ID uuid.UUID) error {
	return svc.delete(ctx, false, ID)
}

// DeleteTx sets deleted_at for a single callhandlings row within a tx from ctx
func (svc *callhandlingService) DeleteTx(ctx context.Context, ID uuid.UUID) error {
	return svc.delete(ctx, true, ID)
}

// delete a callhandling by setting deleted at. if useTx = true then it will attempt to delete the callhandling within a transaction
// from context.
func (svc *callhandlingService) delete(ctx context.Context, useTx bool, ID uuid.UUID) error {
	errMsg := func() string { return "Error executing delete callhandling - " + ID.String() }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["delete-callhandling"])
	} else {
		stmt = svc.stmts["delete-callhandling"]
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
		return errors.Wrap(ErrNotFound, errMsg())
	}

	return nil
}

