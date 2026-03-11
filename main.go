package main

import (
	"database/sql"
	"log"

	"github.com/cristianemek/go-simplebank/api"
	db "github.com/cristianemek/go-simplebank/db/sqlc"
	"github.com/cristianemek/go-simplebank/util"

	_ "github.com/lib/pq" // driver de postgres para conectar con la base de datos
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource) //creamos coexion a la base de datos
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)

}
