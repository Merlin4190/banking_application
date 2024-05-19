package test

import (
	"database/sql"
	"github.com/stretchr/testify/mock"
)

// MockDBContext is a mock implementation of the DBContext interface
type MockDBContext struct {
	mock.Mock
}

type MockRow struct {
	mock.Mock
}

func (r *MockRow) Scan(dest ...interface{}) error {
	args := r.Called(dest)
	for i, d := range dest {
		if args.Get(i) != nil {
			*(d.(*interface{})) = args.Get(i)
		}
	}
	return args.Error(len(dest))
}

// MockTx is a mock implementation of sql.Tx
type MockTx struct {
	mock.Mock
	sql.Tx
}

func (tx *MockTx) Commit() error {
	args := tx.Called()
	return args.Error(0)
}

func (tx *MockTx) Rollback() error {
	args := tx.Called()
	return args.Error(0)
}

type WrappedTx struct {
	*MockTx
	sql.Tx
}

func (m *MockDBContext) Begin() (*sql.Tx, error) {
	args := m.Called()
	tx, _ := args.Get(0).(*sql.Tx)
	return tx, args.Error(1)
}

func (m *MockDBContext) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDBContext) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

// Insert mocks the Insert method of DBContext interface
func (m *MockDBContext) Insert(query string, args ...interface{}) (sql.Result, error) {
	// args[0] contains the query, so we can use it to make assertions if needed
	// Return values expected by the caller
	argsList := make([]interface{}, 0, len(args)+1)
	argsList = append(argsList, query)
	argsList = append(argsList, args...)
	returnValues := m.Called(argsList...)
	return returnValues.Get(0).(sql.Result), returnValues.Error(1)
}

// Update mocks the Update method of DBContext interface
func (m *MockDBContext) Update(query string, args ...interface{}) (sql.Result, error) {
	argsList := make([]interface{}, 0, len(args)+1)
	argsList = append(argsList, query)
	argsList = append(argsList, args...)
	returnValues := m.Called(argsList...)
	return returnValues.Get(0).(sql.Result), returnValues.Error(1)
}

// Query mocks the Query method of DBContext interface
func (m *MockDBContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
	argsList := make([]interface{}, 0, len(args)+1)
	argsList = append(argsList, query)
	argsList = append(argsList, args...)
	returnValues := m.Called(argsList...)
	return returnValues.Get(0).(*sql.Rows), returnValues.Error(1)
}

// QueryRow mocks the QueryRow method of DBContext interface
func (m *MockDBContext) QueryRow(query string, args ...interface{}) *sql.Row {
	argsList := make([]interface{}, 0, len(args)+1)
	argsList = append(argsList, query)
	argsList = append(argsList, args...)
	returnValue := m.Called(argsList...)
	return returnValue.Get(0).(*sql.Row)
}
