package main

import (
	"log"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/volskyi-dmytro/st-bank/api"
	db "github.com/volskyi-dmytro/st-bank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/st_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Can not connect to the database: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
