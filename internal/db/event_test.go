package db

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/caring/call-handling/pb"
	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	event := &pb.Event{
		CallId:     int64(1000),
		EventType:  "dialing",
		IdentityId: int64(2000),
		Timestamp:  int64(20200101),
		Meta:       "twilio_meta",
	}
	proto := &pb.EventRequest{
		Event: event,
	}

	r := NewEvent(proto, "dialing")
	protoEvent := proto.GetEvent()

	assert.Equal(t, protoEvent.GetCallId(), r.CallID, "Expected IDs to match")
	assert.Equal(t, protoEvent.GetEventType(), r.Type, "Expected Types to match")
	assert.Equal(t, protoEvent.GetIdentityId(), r.IdentityID, "Expected Identities to match")
	assert.Equal(t, protoEvent.GetTimestamp(), r.Timestamp, "Expected Timestamps to match")
	assert.Equal(t, protoEvent.GetMeta(), r.Meta, "Expected Meta to match")
}

func TestEvent_create(t *testing.T) {
	stmt := map[string]string{
		"create-event": "INSERT events",
	}
	input := &Event{
		CallID:     int64(2000),
		Type:       "ringing",
		IdentityID: int64(9090),
		Timestamp:  int64(20200101),
		Meta:       "twilio_meta",
	}
	args := []driver.Value{
		int64(2000),
		"ringing",
		int64(9090),
		int64(20200101),
		"twilio_meta",
	}

	// ensures that execution within a transaction occurs without error
	t.Run("With a provided transaction", func(t *testing.T) {
		store, mock, err := NewTestDB(stmt)
		if ok := assert.NoError(t, err, "Expected no error"); !ok {
			assert.FailNow(t, "test setup failed")
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT events").
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(0, 1))

		tx, err := store.GetTx()
		if ok := assert.NoError(t, err, "Expected no error"); !ok {
			assert.FailNow(t, "transaction setup failed")
		}

		err = store.Events.CreateTx(ToCtx(context.Background(), tx), input)
		assert.NoError(t, err, "Expecting no query error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err, "Expecting all mock conditions to be met")
	})

	// ensures that execution outside of a transaction occurs without error
	t.Run("Without a provided transaction", func(t *testing.T) {
		store, mock, err := NewTestDB(stmt)
		if ok := assert.NoError(t, err, "Expected no error"); !ok {
			assert.FailNow(t, "test setup failed")
		}

		mock.ExpectExec("INSERT events").
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = store.Events.Create(context.Background(), input)
		assert.NoError(t, err, "Expecting no query error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err, "Expecting all mock conditions to be met")
	})

	// ensures that a failed record create is handled correctly
	t.Run("Failed record create", func(t *testing.T) {
		store, mock, err := NewTestDB(stmt)
		if ok := assert.NoError(t, err, "Expected no error"); !ok {
			assert.FailNow(t, "test setup failed")
		}

		mock.ExpectExec("INSERT events").
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = store.Events.Create(context.Background(), input)
		assert.EqualError(t, err, "Error executing create call - &{2000, ringing, 9090, 20200101, twilio_meta}: no new rows were created", "Expecting no query error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err, "Expecting all mock conditions to be met")
	})
}
