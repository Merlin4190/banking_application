package database

import "database/sql"

type DBContext interface {
	Begin() (*sql.Tx, error)
	Insert(query string, args ...interface{}) (sql.Result, error)
	Update(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type AppDBContext struct {
	DB *sql.DB
}

func NewDBContext(db *sql.DB) *AppDBContext {
	return &AppDBContext{DB: db}
}

func (ctx *AppDBContext) Begin() (*sql.Tx, error) {
	return ctx.DB.Begin()
}

// Insert executes an insert query and returns the result
func (ctx *AppDBContext) Insert(query string, args ...interface{}) (sql.Result, error) {
	return ctx.DB.Exec(query, args...)
}

// Update executes an update query and returns the result
func (ctx *AppDBContext) Update(query string, args ...interface{}) (sql.Result, error) {
	return ctx.DB.Exec(query, args...)
}

// Query executes a query and returns the result set
func (ctx *AppDBContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return ctx.DB.Query(query, args...)
}

// QueryRow executes a query that is expected to return at most one row
func (ctx *AppDBContext) QueryRow(query string, args ...interface{}) *sql.Row {
	return ctx.DB.QueryRow(query, args...)
}
