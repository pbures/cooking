package tests

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"cooking.buresovi.net/src/persistence"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var dbConnection *sql.DB

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, e error) {

	createDB("testing")

	dbConf := &persistence.Database{
		Username: "postgres",
		Hostname: "localhost",
		Port:     5432,
		Password: "welcome1",
		Database: "testing",
	}

	dbConnection = dbConf.ConnectDb()

	driver, _ := postgres.WithInstance(dbConnection, &postgres.Config{})
	mi, err := migrate.NewWithDatabaseInstance("file://../../database/migrations", "testing", driver)
	if err != nil {
		log.Fatal(err)
	}

	mi.Up()

	defer func() {
		mi.Down()
		dbConnection.Close()
	}()

	return m.Run(), nil
}

func createDB(dbName string) {

	d := &persistence.Database{
		Username: "postgres",
		Hostname: "localhost",
		Port:     5432,
		Password: "welcome1",
	}

	db := d.ConnectDb()
	defer db.Close()

	db.Exec("DROP DATABASE testing")
	smtm := "CREATE DATABASE " + dbName
	_, err := db.Exec(smtm)
	if err != nil {
		log.Fatal(err)
	}
}
