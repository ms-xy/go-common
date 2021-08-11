package transactions

import (
  "context"
  "database/sql"
  "database/sql/driver"
  "errors"
  "testing"

  "github.com/ms-xy/go-common/database/mockdb"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
)

func any(s string) mock.AnythingOfTypeArgument {
  return mock.AnythingOfType(s)
}

func setupTx(tx *mockdb.Tx, commitErr, rollbackErr error) {
  *tx = mockdb.Tx{}
  tx.On("Commit").Return(commitErr)
  tx.On("Rollback").Return(rollbackErr)
}

func TestWithTx(t *testing.T) {
  // Test Setup ------------------------------------------------------------------------------------
  mockDriver := new(mockdb.Driver)
  sql.Register("testdriver", mockDriver)
  db, err := sql.Open("testdriver", "")
  if err != nil {
    panic(err)
  }
  defer db.Close()

  // MOCKS -----------------------------------------------------------------------------------------
  conn := new(mockdb.Conn)
  tx := new(mockdb.Tx) // will be reset anyways, but it's necessary to assign an address

  // TMP Test Vars ---------------------------------------------------------------------------------
  var tmpCtx context.Context
  var tmpOpts driver.TxOptions
  var tmpFnExecuted bool = false

  // MOCK BEHAVIOR ---------------------------------------------------------------------------------
  mockDriver.
    On("Open", any("string")).
    Return(conn, nil)
  conn.
    On("BeginTx", any("*context.cancelCtx"), any("driver.TxOptions")).
    Run(func(args mock.Arguments) {
      tmpCtx = args.Get(0).(context.Context)
      tmpOpts = args.Get(1).(driver.TxOptions)
    }).
    Return(tx, nil).
    On("Close").
    Return(nil)

  // EXECUTE TESTS ---------------------------------------------------------------------------------

  setupTx(tx, nil, nil) // reset tx
  WithReadTx(db, nil, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
    // When calling WithTx then a valid context must be passed onto fn
    require.NotNil(t, ctx)
    // When calling WithTx then valid tx should be passed to fn
    require.NotNil(t, tx)
    // indicate fn executed
    tmpFnExecuted = true
    return nil, nil
  })
  // When calling WithTx then db.BeginTx should be called
  conn.AssertNumberOfCalls(t, "BeginTx", 1)
  // When calling WithTx and ctx is nil, a new context should be created
  require.NotNil(t, tmpCtx)
  // When calling WithTx with opts, then opts should be passed on to db.BeginTx
  require.NotNil(t, tmpOpts)
  require.Equal(t, driver.IsolationLevel(sql.LevelDefault), tmpOpts.Isolation)
  require.Equal(t, true, tmpOpts.ReadOnly)
  // When calling WithTx then fn should be executed
  require.True(t, tmpFnExecuted)
  // When calling WithTx and fn succeeds then tx should be committed
  tx.AssertNumberOfCalls(t, "Commit", 1)
  tx.AssertNumberOfCalls(t, "Rollback", 0)

  // When calling WithTx and fn fails with an error then tx should be rolled back
  setupTx(tx, nil, nil) // reset tx
  _, err = WithWriteTx(db, nil, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
    return nil, errors.New("fn error")
  })
  tx.AssertNumberOfCalls(t, "Rollback", 1)
  tx.AssertNumberOfCalls(t, "Commit", 0)
  require.NotNil(t, err)
  require.Equal(t, "fn error", err.Error())

  // When calling WithTx and fn panics then tx should be rolled back and panic rethrown
  setupTx(tx, nil, nil) // reset tx
  require.PanicsWithError(t, "PANIC", func() {
    WithWriteTx(db, nil, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
      panic("PANIC")
    })
  })
  tx.AssertNumberOfCalls(t, "Rollback", 1) // second test that requires a rollback
  tx.AssertNumberOfCalls(t, "Commit", 0)

  // When calling WithTx and fn panics and rollback fails then there should be a wrapped error
  setupTx(tx, nil, errors.New("rollback error"))
  require.PanicsWithError(t, "rollback error: PANIC", func() {
    WithWriteTx(db, nil, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
      panic("PANIC")
    })
  })
  tx.AssertNumberOfCalls(t, "Rollback", 1) // second test that requires a rollback
  tx.AssertNumberOfCalls(t, "Commit", 0)

  // When calling WithTx and commit fails then the commit error should be returned
  setupTx(tx, errors.New("commit error"), nil) // reset tx
  _, err = WithReadTx(db, nil, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
    return nil, nil
  })
  tx.AssertNumberOfCalls(t, "Commit", 1)
  tx.AssertNumberOfCalls(t, "Rollback", 0)
  require.NotNil(t, err)
  require.Equal(t, "commit error", err.Error())

  // When calling WithTx and an error happens and rollback fails it should wrap the second error
  setupTx(tx, nil, errors.New("rollback error")) // reset tx
  _, err = WithReadTx(db, nil, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
    return nil, errors.New("fn error")
  })
  tx.AssertNumberOfCalls(t, "Commit", 0)
  tx.AssertNumberOfCalls(t, "Rollback", 1)
  require.NotNil(t, err)
  require.Equal(t, "rollback error: fn error", err.Error())
}
