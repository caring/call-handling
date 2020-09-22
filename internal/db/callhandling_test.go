package db

import (
  "context"
  "database/sql"
  "database/sql/driver"
  "testing"

  "github.com/DATA-DOG/go-sqlmock"
  "github.com/google/uuid"
  "github.com/stretchr/testify/assert"

  "github.com/caring/call-handling/pb"
)



// ensures that casting from proto to store structs occurs correctly
func TestNewCallhandling(t *testing.T) {
  callhandlingID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
  proto := pb.CreateCallhandlingRequest{
    Name:       "Foobar",
  }

  r, err := NewCallhandling(callhandlingID.String(), &proto)

  assert.NoError(t, err, "Expected NewCategory not to error")
  assert.Equal(t, callhandlingID, r.ID, "Expected UUIDs to match")
  assert.Equal(t, proto.Name, r.Name, "Expected name to be correctly assigned")
}

// ensures that casting from store to proto response occurs correctly
func TestCallhandling_ToProto(t *testing.T) {
  callhandlingID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")

  callhandling := &Callhandling{
    ID:  callhandlingID,
    Name:       "foobar",
  }

  r := callhandling.ToProto()

  assert.Equal(t, callhandlingID.String(), r.CallhandlingId, "Expected field to be mapped back to proto object correctly")
  assert.Equal(t, "foobar", r.Name, "Expected field to be mapped back to proto object correctly")
}

func TestCallhandlingService_get(t *testing.T) {
  callhandlingID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
  stmt := map[string]string{
    "get-callhandling": "SELECT callhandlings",
  }
  args := []driver.Value{
    "72bc87f3-4a9f-4d05-93fe-844d3cd94c65",
  }

  // ensures execution within a transaction occurs without error and the correct result is returned
  t.Run("With a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectBegin()
    mock.ExpectQuery("SELECT callhandlings").
      WithArgs(args...).
      WillReturnRows(
        sqlmock.NewRows([]string{"callhandling_id", "name"}).
          AddRow(callhandlingID, "Foobar"),
      )

    tx, err := store.GetTx()
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "transaction setup failed")
    }

    result, err := store.Callhandling.GetTx(ToCtx(context.Background(), tx), callhandlingID)
    assert.NoError(t, err, "Expecting no query error")

    assert.Equal(t, callhandlingID, r.ID, "Expected correct callhandling ID to be returned")
    assert.Equal(t, "Foobar", r.Name, "Expected correct name to be returned")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures that execution outside of transaction occurs without error and the correct result is returned
  t.Run("Without a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectQuery("SELECT callhandlings").
      WithArgs(args...).
      WillReturnRows(
        sqlmock.NewRows([]string{"callhandling_id", "name"}).
          AddRow(callhandlingID, "Foobar"),
      )

    result, err := store.Callhandling.Get(context.Background(), callhandlingID)
    assert.NoError(t, err, "Expecting no query error")

    assert.Equal(t, callhandlingID, r.ID, "Expected correct callhandling ID to be returned")
    assert.Equal(t, "Foobar", r.Name, "Expected correct name to be returned")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures a record not found is handled correctly
  t.Run("No rows returned", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectQuery("SELECT callhandlings").
      WithArgs(args...).WillReturnError(sql.ErrNoRows)

    _, err = store.Callhandling.Get(context.Background(), callhandlingID)
    assert.EqualError(t, err, "Error executing get callhandling - 72bc87f3-4a9f-4d05-93fe-844d3cd94c65: the record you are attempting to find or update is not found", "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })
}

func TestCallhandlingService_create(t *testing.T) {
  callhandlingID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
  stmt := map[string]string{
    "create-callhandling": "INSERT callhandlings",
  }
  input := &Callhandling{
    ID:   callhandlingID,
    Name: "Foobar",
  }
  args := []driver.Value{
    "72bc87f3-4a9f-4d05-93fe-844d3cd94c65",
    "Foobar",
  }

  // ensures that execution within a transaction occurs without error
  t.Run("With a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectBegin()
    mock.ExpectExec("INSERT callhandlings").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    tx, err := store.GetTx()
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "transaction setup failed")
    }

    err = store.Callhandling.CreateTx(ToCtx(context.Background(), tx), input)
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

    mock.ExpectExec("INSERT callhandlings").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    err = store.Callhandling.Create(context.Background(), input)
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

    mock.ExpectExec("INSERT callhandlings").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 0))

    err = store.Callhandling.Create(context.Background(), input)
    assert.EqualError(t, err, "Error executing create callhandling - &{72bc87f3-4a9f-4d05-93fe-844d3cd94c65 Foobar}: no new rows were created", "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })
}

func TestCallhandlingService_update(t *testing.T) {
  callhandlingID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
  stmt := map[string]string{
    "update-callhandling": "UPDATE callhandlings",
  }
  input := &Callhandling{
    ID:   callhandlingID,
    Name: "Foobar",
  }
  args := []driver.Value{
    "Foobar",
    "72bc87f3-4a9f-4d05-93fe-844d3cd94c65",
  }

  // ensures that execution within a transaction occurs without error
  t.Run("With a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectBegin()
    mock.ExpectExec("UPDATE callhandlings").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    tx, err := store.GetTx()
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "transaction setup failed")
    }

    err = store.Callhandling.UpdateTx(ToCtx(context.Background(), tx), input)
    assert.NoError(t, err, "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures execution out of a transaction occurs without error
  t.Run("Without a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectExec("UPDATE callhandlings").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    err = store.Callhandling.Update(context.Background(), input)
    assert.NoError(t, err, "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures correct error to be returned when no rows are updated
  t.Run("No updates occurred", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectExec("UPDATE callhandlings").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 0))

    err = store.Callhandling.Update(context.Background(), input)
    assert.EqualError(t, err, "Error executing update callhandling - &{72bc87f3-4a9f-4d05-93fe-844d3cd94c65 Foobar}: no rows affected", "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })
}

func TestCallhandlingService_delete(t *testing.T) {
  callhandlingID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
  stmt := map[string]string{
    "delete-callhandling": "UPDATE callhandlings",
  }
  args := []driver.Value{
    "72bc87f3-4a9f-4d05-93fe-844d3cd94c65",
  }

  // ensures that execution withing a transaction occurs without error
  t.Run("With a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectBegin()
    mock.ExpectExec("UPDATE callhandlings").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    tx, err := store.GetTx()
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "transaction setup failed")
    }

    err = store.Callhandling.DeleteTx(ToCtx(context.Background(), tx), callhandlingID)
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

    mock.ExpectExec("UPDATE callhandlings").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    err = store.Callhandling.Delete(context.Background(), callhandlingID)
    assert.NoError(t, err, "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures that deleting a non existent record is handled correctly
  t.Run("Deleting a non existent record", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectExec("UPDATE callhandlings").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 0))

    err = store.Callhandling.Delete(context.Background(), callhandlingID)
    assert.EqualError(t, err, "Error executing delete callhandling - 72bc87f3-4a9f-4d05-93fe-844d3cd94c65: the record you are attempting to find or update is not found", "Expecting not found error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })
}