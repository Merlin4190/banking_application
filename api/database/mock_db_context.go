package database

import (
	"database/sql"
	"github.com/stretchr/testify/mock"
)

// MockDBContext is a mock implementation of the DBContext interface
type MockDBContext struct {
	mock.Mock
}

func (m *MockDBContext) Begin() (*sql.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sql.Tx), args.Error(1)
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
