package tests

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/persistence"
	"cooking.buresovi.net/src/persistence/meal"
	"cooking.buresovi.net/src/persistence/user"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx"
)

var dbConnection *sql.DB
var application app.App

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, e error) {

	dbConf := &persistence.Database{
		Username: "postgres",
		Hostname: "localhost",
		Port:     5432,
		Password: "welcome1",
	}

	persistence.CreateDB(*dbConf, "testing")
	dbConf.Database = "testing"

	dbConnection = dbConf.ConnectDb()

	driver, _ := postgres.WithInstance(dbConnection, &postgres.Config{})
	mi, err := migrate.NewWithDatabaseInstance("file://../../database/migrations", "testing", driver)
	if err != nil {
		log.Fatal(err)
	}

	mi.Up()

	defer func() {
		// mi.Down()
		dbConnection.Close()
	}()

	userSvc := user.NewUserSvcPsql(dbConnection)
	mealsSvc := meal.NewMealSvcPsql(userSvc, dbConnection)

	application = app.App{
		MealSvc: mealsSvc,
		UserSvc: userSvc,
	}

	return m.Run(), nil
}
