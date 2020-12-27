package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var (
	// SQL wrapper
	SQL *sql.DB
	// Database info
	databases PostGresSQL
)

// Type is the type of database from a Type* constant
type Type string

const (
	TypePostGresSQL Type = "PostGresSQL"
)

// Info contains the database configurations
type Info struct {
	PostGresSQL PostGresSQL
}

// MySQLInfo is the details for the database connection
type PostGresSQL struct {
	Username string
	Password string
	DBName   string
	Hostname string
	SSLMode  bool
}

// DSN returns the Data Source Name
func DSN(ci PostGresSQL) string {
	// Example: postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full
	return "postgres://" +
		ci.Username +
		":" +
		ci.Password +
		"@" +
		ci.Hostname +
		"/" +
		ci.DBName +
		"?sslmode=disable"
}

// Connect to the database
func Connect(d PostGresSQL) {
	var err error

	// Store the config
	databases = d

	// Connect to Postgres
	if SQL, err = sql.Open("postgres", DSN(d)); err != nil {
		log.Println("SQL Driver Error", err)
	}

	// Check if is alive
	if err = SQL.Ping(); err != nil {
		log.Println("Database Error", err)
	}
}
