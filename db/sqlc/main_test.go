package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	dbDriver   string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     int
	dbName     string
	dbSource   string
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func init() {
	var err error
	godotenv.Load("../../.env")

	dbDriver = os.Getenv("DB_DRIVER")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbHost = os.Getenv("DB_HOST")
	dbPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("invalid DB port: %v", err)
	}

	dbName = os.Getenv("DB_NAME")
	dbSource = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", dbDriver, dbUser, dbPassword, dbHost, dbPort, dbName)
}

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
