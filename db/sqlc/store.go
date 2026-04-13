package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store que contiene todas las funciones para ejecutar consultas a la base de datos y transacciones
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

// SQLStore proporciona todas las funciones para ejecutar consultas SQL y transacciones
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// ejecutar una funcion con una database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err. %w, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
