package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/caring/call-handling/pb"
)

// ensures that casting from proto to store structs occurs correctly
func TestNewCall(t *testing.T) {
	call := &pb.Call{
		CallId:         int64(1000),
		Sid:            int64(2000),
		ConversationId: int64(3000),
		ANI:            "1011011111",
		DNIS:           "9099099999",
		Status:         "active",
	}
	proto := &pb.CreateCallRequest{
		Call: call,
	}

	r := NewCall(proto)

	protoCall := proto.GetCall()

	assert.Equal(t, protoCall.GetCallId(), r.ID, "Expected IDs to match")
	assert.Equal(t, protoCall.GetSid(), r.SID, "Expected SID to match")
	assert.Equal(t, protoCall.GetConversationId(), r.ConversationID, "Expected ConversationId to match")
	assert.Equal(t, protoCall.GetANI(), r.ANI, "Expected ANI to match")
	assert.Equal(t, protoCall.GetDNIS(), r.DNIS, "Expected DNIS to match")
	assert.Equal(t, protoCall.GetStatus(), r.Status, "Expected Status to match")
}

// ensures that casting from store to proto response occurs correctly
// func TestCallhandling_ToProto(t *testing.T) {
// 	call := &Call{
// 		ID:         int64(1000),
// 	}

// 	r := callhandling.ToProto()

// 	assert.Equal(t, callhandlingID.String(), r.CallhandlingId, "Expected field to be mapped back to proto object correctly")
// }

func TestCallhandlingService_get(t *testing.T) {
	callID := int64(2000)
	sid := int64(3000)
	conv_id := int64(4000)
	ani := "2022022222"
	dnis := "8088088888"
	status := "active"
	stmt := map[string]string{
		"get-call": "SELECT calls",
	}
	args := []driver.Value{
		int64(2000),
		int64(3000),
		int64(4000),
		"2022022222",
		"8088088888",
		"active",
	}

	// ensures execution within a transaction occurs without error and the correct result is returned
	t.Run("With a provided transaction", func(t *testing.T) {
		store, mock, err := NewTestDB(stmt)
		if ok := assert.NoError(t, err, "Expected no error"); !ok {
			assert.FailNow(t, "test setup failed")
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT calls").
			WithArgs(args...).
			WillReturnRows(
				sqlmock.NewRows([]string{"call_id", "sid", "conversation_id", "ANI", "DNIS", "status"}).
					AddRow(callID, sid, conv_id, ani, dnis, status),
			)

		tx, err := store.GetTx()
		if ok := assert.NoError(t, err, "Expected no error"); !ok {
			assert.FailNow(t, "transaction setup failed")
		}

		r, err := store.Calls.GetTx(ToCtx(context.Background(), tx), callID)
		assert.NoError(t, err, "Expecting no query error")

		assert.Equal(t, callID, r.ID, "Expected correct call ID to be returned")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err, "Expecting all mock conditions to be met")
	})

	// ensures that execution outside of transaction occurs without error and the correct result is returned
	t.Run("Without a provided transaction", func(t *testing.T) {
		store, mock, err := NewTestDB(stmt)
		if ok := assert.NoError(t, err, "Expected no error"); !ok {
			assert.FailNow(t, "test setup failed")
		}

		mock.ExpectQuery("SELECT calls").
			WithArgs(args...).
			WillReturnRows(
				sqlmock.NewRows([]string{"call_id", "sid", "conversation_id", "ANI", "DNIS", "status"}).
					AddRow(callID, sid, conv_id, ani, dnis, status),
			)

		r, err := store.Calls.Get(context.Background(), callID)
		assert.NoError(t, err, "Expecting no query error")

		assert.Equal(t, callID, r.ID, "Expected correct callhandling ID to be returned")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err, "Expecting all mock conditions to be met")
	})

	// ensures a record not found is handled correctly
	t.Run("No rows returned", func(t *testing.T) {
		store, mock, err := NewTestDB(stmt)
		if ok := assert.NoError(t, err, "Expected no error"); !ok {
			assert.FailNow(t, "test setup failed")
		}

		mock.ExpectQuery("SELECT calls").
			WithArgs(args...).WillReturnError(sql.ErrNoRows)

		_, err = store.Calls.Get(context.Background(), callID)
		assert.EqualError(t, err, "Error executing get call - 2000: the record you are attempting to find or update is not found", "Expecting no query error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err, "Expecting all mock conditions to be met")
	})
}

func TestCallhandlingService_create(t *testing.T) {
	stmt := map[string]string{
		"create-call": "INSERT calls",
	}
	input := &Call{
		ID:             int64(3000),
		SID:            int64(4000),
		ConversationID: int64(5000),
		ANI:            "3033033333",
		DNIS:           "7077077777",
		Status:         "active",
	}
	args := []driver.Value{
		int64(3000),
		int64(4000),
		int64(5000),
		"3033033333",
		"7077077777",
		"active",
	}

	// ensures that execution within a transaction occurs without error
	t.Run("With a provided transaction", func(t *testing.T) {
		store, mock, err := NewTestDB(stmt)
		if ok := assert.NoError(t, err, "Expected no error"); !ok {
			assert.FailNow(t, "test setup failed")
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT calls").
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(0, 1))

		tx, err := store.GetTx()
		if ok := assert.NoError(t, err, "Expected no error"); !ok {
			assert.FailNow(t, "transaction setup failed")
		}

		err = store.Calls.CreateTx(ToCtx(context.Background(), tx), input)
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

		mock.ExpectExec("INSERT calls").
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = store.Calls.Create(context.Background(), input)
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

		mock.ExpectExec("INSERT calls").
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = store.Calls.Create(context.Background(), input)
		assert.EqualError(t, err, "Error executing create call - &{3000, 4000, 5000, 3033033333, 7077077777, active}: no new rows were created", "Expecting no query error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err, "Expecting all mock conditions to be met")
	})
}

// func TestCallhandlingService_update(t *testing.T) {
//   callhandlingID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
//   stmt := map[string]string{
//     "update-callhandling": "UPDATE callhandlings",
//   }
//   input := &Callhandling{
//     ID:   callhandlingID,
//     Name: "Foobar",
//   }
//   args := []driver.Value{
//     "Foobar",
//     "72bc87f3-4a9f-4d05-93fe-844d3cd94c65",
//   }

//   // ensures that execution within a transaction occurs without error
//   t.Run("With a provided transaction", func(t *testing.T) {
//     store, mock, err := NewTestDB(stmt)
//     if ok := assert.NoError(t, err, "Expected no error"); !ok {
//       assert.FailNow(t, "test setup failed")
//     }

//     mock.ExpectBegin()
//     mock.ExpectExec("UPDATE callhandlings").
//       WithArgs(args...).
//       WillReturnResult(sqlmock.NewResult(0, 1))

//     tx, err := store.GetTx()
//     if ok := assert.NoError(t, err, "Expected no error"); !ok {
//       assert.FailNow(t, "transaction setup failed")
//     }

//     err = store.Callhandling.UpdateTx(ToCtx(context.Background(), tx), input)
//     assert.NoError(t, err, "Expecting no query error")

//     err = mock.ExpectationsWereMet()
//     assert.NoError(t, err, "Expecting all mock conditions to be met")
//   })

//   // ensures execution out of a transaction occurs without error
//   t.Run("Without a provided transaction", func(t *testing.T) {
//     store, mock, err := NewTestDB(stmt)
//     if ok := assert.NoError(t, err, "Expected no error"); !ok {
//       assert.FailNow(t, "test setup failed")
//     }

//     mock.ExpectExec("UPDATE callhandlings").
//       WithArgs(args...).
//       WillReturnResult(sqlmock.NewResult(0, 1))

//     err = store.Callhandling.Update(context.Background(), input)
//     assert.NoError(t, err, "Expecting no query error")

//     err = mock.ExpectationsWereMet()
//     assert.NoError(t, err, "Expecting all mock conditions to be met")
//   })

//   // ensures correct error to be returned when no rows are updated
//   t.Run("No updates occurred", func(t *testing.T) {
//     store, mock, err := NewTestDB(stmt)
//     if ok := assert.NoError(t, err, "Expected no error"); !ok {
//       assert.FailNow(t, "test setup failed")
//     }

//     mock.ExpectExec("UPDATE callhandlings").
//       WithArgs(args...).
//       WillReturnResult(sqlmock.NewResult(0, 0))

//     err = store.Callhandling.Update(context.Background(), input)
//     assert.EqualError(t, err, "Error executing update callhandling - &{72bc87f3-4a9f-4d05-93fe-844d3cd94c65 Foobar}: no rows affected", "Expecting no query error")

//     err = mock.ExpectationsWereMet()
//     assert.NoError(t, err, "Expecting all mock conditions to be met")
//   })
// }

// func TestCallhandlingService_delete(t *testing.T) {
//   callhandlingID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
//   stmt := map[string]string{
//     "delete-callhandling": "UPDATE callhandlings",
//   }
//   args := []driver.Value{
//     "72bc87f3-4a9f-4d05-93fe-844d3cd94c65",
//   }

//   // ensures that execution withing a transaction occurs without error
//   t.Run("With a provided transaction", func(t *testing.T) {
//     store, mock, err := NewTestDB(stmt)
//     if ok := assert.NoError(t, err, "Expected no error"); !ok {
//       assert.FailNow(t, "test setup failed")
//     }

//     mock.ExpectBegin()
//     mock.ExpectExec("UPDATE callhandlings").
//       WithArgs(args...).
//       WillReturnResult(sqlmock.NewResult(0, 1))

//     tx, err := store.GetTx()
//     if ok := assert.NoError(t, err, "Expected no error"); !ok {
//       assert.FailNow(t, "transaction setup failed")
//     }

//     err = store.Callhandling.DeleteTx(ToCtx(context.Background(), tx), callhandlingID)
//     assert.NoError(t, err, "Expecting no query error")

//     err = mock.ExpectationsWereMet()
//     assert.NoError(t, err, "Expecting all mock conditions to be met")
//   })

//   // ensures that execution outside of a transaction occurs without error
//   t.Run("Without a provided transaction", func(t *testing.T) {
//     store, mock, err := NewTestDB(stmt)
//     if ok := assert.NoError(t, err, "Expected no error"); !ok {
//       assert.FailNow(t, "test setup failed")
//     }

//     mock.ExpectExec("UPDATE callhandlings").
//       WithArgs(args...).
//       WillReturnResult(sqlmock.NewResult(0, 1))

//     err = store.Callhandling.Delete(context.Background(), callhandlingID)
//     assert.NoError(t, err, "Expecting no query error")

//     err = mock.ExpectationsWereMet()
//     assert.NoError(t, err, "Expecting all mock conditions to be met")
//   })

//   // ensures that deleting a non existent record is handled correctly
//   t.Run("Deleting a non existent record", func(t *testing.T) {
//     store, mock, err := NewTestDB(stmt)
//     if ok := assert.NoError(t, err, "Expected no error"); !ok {
//       assert.FailNow(t, "test setup failed")
//     }

//     mock.ExpectExec("UPDATE callhandlings").
//       WithArgs(args...).
//       WillReturnResult(sqlmock.NewResult(0, 0))

//     err = store.Callhandling.Delete(context.Background(), callhandlingID)
//     assert.EqualError(t, err, "Error executing delete callhandling - 72bc87f3-4a9f-4d05-93fe-844d3cd94c65: the record you are attempting to find or update is not found", "Expecting not found error")

//     err = mock.ExpectationsWereMet()
//     assert.NoError(t, err, "Expecting all mock conditions to be met")
//   })
// }
