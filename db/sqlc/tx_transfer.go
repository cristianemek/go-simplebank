package db

import (
	"context"
)

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

// var txKey = struct{}{} // clave para almacenar la transaccion en el contexto creamos objeto vacio

// TrnsferTx para realizar una transferencia de dinero entre dos cuentas.
// Esta función realiza varias operaciones dentro de una transacción: crear una transferencia, crear entradas para ambas cuentas y actualizar los balances de las cuentas. Si alguna de estas operaciones falla, la transacción se revertirá para mantener la integridad de los datos.
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult //creamos resultado vacio

	err := store.execTx(ctx, func(q *Queries) error { //este objeto de consulta se crea a partir de la transacción, por lo que todas las operaciones dentro de esta función se ejecutarán dentro de la misma transacción
		var err error

		// txName := ctx.Value(txKey) //obtenemos el nombre de la transaccion del contexto para identificarla en los logs

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

		if arg.FromAccountID < arg.ToAccountID { // para evitar deadlocks, siempre actualizamos las cuentas en el mismo orden, en este caso, de menor a mayor ID
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
