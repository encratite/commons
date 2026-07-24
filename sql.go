package commons

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SQLConnection struct {
	db *sql.DB
}

type SQLRows struct {
	rows *sql.Rows
}

func NewSQLConnection(dataSourceName string) SQLConnection {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		Fatalf("Failed to connect: %v", err)
	}
	return SQLConnection{
		db: db,
	}
}

func (c *SQLConnection) Close() {
	c.db.Close()
}

func (c *SQLConnection) Query(query string, args ...any) SQLRows {
	rows, err := c.db.Query(query, args...)
	if err != nil {
		Fatalf("Failed to perform query: %v", err)
	}
	return SQLRows{
		rows: rows,
	}
}

func (r *SQLRows) Read(onRead func (), args ...any) time.Duration {
	start := time.Now()
	for r.rows.Next() {
		err := r.rows.Scan(args...)
		onRead()
		if err != nil {
			Fatalf("Failed to read rows: %v", err)
		}
	}
	err := r.rows.Err()
	if err != nil {
		Fatalf("Failed to read rows: %v", err)
	}
	end := time.Now()
	duration := end.Sub(start)
	return duration
}