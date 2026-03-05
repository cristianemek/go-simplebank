package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store para combinar todas las consultas generadas por sqlc y agregar funciones personalizadas como transacciones
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// ejecutar una funcion con una database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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

// TrnasferTxParams contiene los parametros necesarios para realizar una transferencia de dinero entre dos cuentas
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contiene los resultados de una transferencia de dinero entre dos cuentas
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
}

// TrnsferTx para realizar una transferencia de dinero entre dos cuentas.
// Esta función realiza varias operaciones dentro de una transacción: crear una transferencia, crear entradas para ambas cuentas y actualizar los balances de las cuentas. Si alguna de estas operaciones falla, la transacción se revertirá para mantener la integridad de los datos.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult //creamos resultado vacio

	err := store.execTx(ctx, func(q *Queries) error { //este objeto de consulta se crea a partir de la transacción, por lo que todas las operaciones dentro de esta función se ejecutarán dentro de la misma transacción
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{ //creamos la transferencia
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{ //creamos la entrada para la cuenta de origen
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{ //creamos la entrada para la cuenta de destino
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//todo actualizar los balances de las cuentas

		return nil
	})
	return result, err
}
