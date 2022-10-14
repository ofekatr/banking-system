package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ofekatr/simple-bank/util"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	var (
		dbDriver = config.DBDriver
		dbSource = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", config.DBDriver, config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	)

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
