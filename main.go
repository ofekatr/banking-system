package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/ofekatr/simple-bank/api"
	db "github.com/ofekatr/simple-bank/db/sqlc"
	"github.com/ofekatr/simple-bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	var (
		serverAddress = config.ServerAddress
		dbDriver      = config.DBDriver
		dbSource      = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", config.DBDriver, config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	)

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
