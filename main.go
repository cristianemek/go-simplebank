package main

import (
	"database/sql"
	"log"

	"github.com/cristianemek/go-simplebank/api"
	db "github.com/cristianemek/go-simplebank/db/sqlc"

	_ "github.com/lib/pq" // driver de postgres para conectar con la base de datos
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource) //creamos coexion a la base de datos
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)

}
