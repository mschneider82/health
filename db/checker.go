package db

import (
	"database/sql"

	"github.com/mschneider82/health"
)

// Checker is a checker that check a database connection
type Checker struct {
	CheckSQL   string
	VersionSQL string
	DB         *sql.DB
}

// NewChecker returns a new db.Checker with the given URL
func NewChecker(checkSQL, versionSQL string, db *sql.DB) Checker {
	return Checker{CheckSQL: checkSQL, VersionSQL: versionSQL, DB: db}
}

// NewMySQLChecker returns a new db.Checker configured for use in MySQL
func NewMySQLChecker(db *sql.DB) Checker {
	return NewChecker("SELECT 1", "SELECT VERSION()", db)
}

// NewPostgreSQLChecker returns a new db.Checker configured for use in PostgreSQL
func NewPostgreSQLChecker(db *sql.DB) Checker {
	return NewChecker("SELECT 1", "SELECT VERSION()", db)
}

// NewSqlite3Checker returns a new db.Checker configured for use in Sqlite3
func NewSqlite3Checker(db *sql.DB) Checker {
	return NewChecker("SELECT 1", "SELECT sqlite_version()", db)
}

// Check execute two queries in the database
// The first is a simple one used to verify if the database is up
// If is Up then another query is executed, querying for the database version
func (c Checker) Check() health.Health {
	var (
		version string
		ok      string
	)

	h := health.NewHealth()

	if c.DB == nil {
		h.Down().AddInfo("error", "Empty resource")
		return h
	}

	err := c.DB.QueryRow(c.CheckSQL).Scan(&ok)

	if err != nil {
		h.Down().AddInfo("error", err.Error())
		return h
	}

	// We are gonna make it VersionSQL optional, as I cannot change the API,
	// we decided to ignore if the VersionSQL is empty.
	if c.VersionSQL != "" {
		err = c.DB.QueryRow(c.VersionSQL).Scan(&version)

		if err != nil {
			h.Down().AddInfo("error", err.Error())
			return h
		}

		h.Up().AddInfo("version", version)
	}

	return h
}
