package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cristianemek/go-simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) { // por convencion la funcion TestMain es el punto de entrada para todas las pruebas unitarias
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource) //creamos coexion a la base de datos
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
