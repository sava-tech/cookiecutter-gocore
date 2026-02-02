package database

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"{{ cookiecutter.module_path }}/utils"
)

var (
	testQueries *Queries
	testDB      *pgxpool.Pool
)

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	defer testDB.Close()
	testQueries = New(testDB)
	// Run the tests and exit appropriately
	os.Exit(m.Run())
}
