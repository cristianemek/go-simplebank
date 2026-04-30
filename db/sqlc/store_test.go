package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf(">> antes: account1: %v, account2: %v\n", account1.Balance, account2.Balance)

	// vamos a realizar 5 transferencias concurrentes de account1 a account2

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ { // estas rutinas se ejecutan en una diferente a TestTransferTx por lo que, el test no hay garantia de que se detenga si una condicion de error se cumple
		// txName := fmt.Sprintf("tx %d", i+1) // nombre de la transaccion para identificarla en los logs
		go func() { // cada transferencia se ejecuta en una goroutine diferente para simular la concurrencia
			// ctx := context.WithValue(context.Background(), txKey, txName) // creamos un contexto con el nombre de la transaccion para identificarla en los logs
			ctx := context.Background() // creamos un contexto
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err       // enviamos el error al canal de errores
			results <- result // enviamos el resultado al canal de resultados
		}()
	}

	existed := make(map[int]bool) // mapa para verificar que cada transferencia es unica
	for range n {
		err := <-errs // recibimos el error del canal de errores
		require.NoError(t, err)

		result := <-results // recibimos el resultado del canal de resultados
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// comprobar las entradas
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//revisar cuentas

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//comprobar balance de las cuentas
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // el monto de la transferencia debe ser un multiplo del monto de cada transferencia

		k := int(diff1 / amount) // el numero de transferencias realizadas hasta ahora

		require.True(t, k >= 1 && k <= n)  // el numero de transferencias realizadas debe ser entre 1 y n
		require.NotContains(t, existed, k) // cada transferencia debe ser unica
		existed[k] = true
	}

	// comprobar el balance final de las cuentas
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">> despues: account1: %v, account2: %v\n", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf(">> antes: account1: %v, account2: %v\n", account1.Balance, account2.Balance)

	// vamos a realizar 5 transferencias concurrentes de account1 a account2

	n := 10
	amount := int64(10)

	errs := make(chan error)

	// para simular un posible deadlock, vamos a realizar 10 transferencias concurrentes en ambas direcciones
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 { // para las transferencias impares, invertimos el orden de las cuentas para simular un posible deadlock
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			ctx := context.Background() // creamos un contexto
			_, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for range n {
		err := <-errs // recibimos el error del canal de errores
		require.NoError(t, err)

	}

	// comprobar el balance final de las cuentas
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">> despues: account1: %v, account2: %v\n", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
