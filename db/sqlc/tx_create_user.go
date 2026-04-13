package db

import (
	"context"
)

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error { //este objeto de consulta se crea a partir de la transacción, por lo que todas las operaciones dentro de esta función se ejecutarán dentro de la misma transacción
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)

		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})
	return result, err
}
