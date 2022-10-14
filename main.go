package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ofekatr/simple-bank/api"
	db "github.com/ofekatr/simple-bank/db/sqlc"
)

var (
	serverAddress string

	dbDriver   string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     int
	dbName     string
	dbSource   string
)

func init() {
	var err error
	godotenv.Load(".env")

	serverAddress = os.Getenv("SERVER_ADDRESS")

	dbDriver = "postgres"
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

func main() {
	var err error

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(serverAddress); err != nil {
		log.Fatalf("cannot start server: %v", err)
	}

	log.Printf("server is listening at %s", serverAddress)
}
