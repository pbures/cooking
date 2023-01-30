package persistence

import (
	"database/sql"
	"fmt"
	"log"
)

type Database struct {
	Hostname string
	Port     int
	Database string
	Username string
	Password string
}

func (d *Database) ConnectDb() *sql.DB {

	connStr := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v sslmode=disable",
		d.Hostname, d.Port, d.Username, d.Password)

	if d.Database != "" {
		connStr = connStr + " database=" + d.Database
	}

	// var err error = nil

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
