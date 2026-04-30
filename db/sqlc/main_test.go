package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cristianemek/go-simplebank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) { // por convencion la funcion TestMain es el punto de entrada para todas las pruebas unitarias
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource) //creamos coexion a la base de datos
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
