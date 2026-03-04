package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// todo por el momento las declaro como constantes en el futuro se deben de leer de las variables de entorno
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) { // por convencion la funcion TestMain es el punto de entrada para todas las pruebas unitarias
	conn, err := sql.Open(dbDriver, dbSource) //creamos coexion a la base de datos
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
