package db

import (
	"context"
	"database/sql"
	"fmt"
)


type Store interface{
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResults, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResults, error) 
	Querier

}

// the store will provide all functions to execute the database queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries

}

// A constructor that create a new instance of the store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

// execTx exuctes a function within a database transaction
func (store *SQLStore) execTX(ctx context.Context, fn func(*Queries) error) error {

	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	
	
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error during transaction %v, %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()

}

type TransferTxParams struct {
	FromAccountId int32 `json:"from_account_id"`
	ToAccountId int32 `json:"to_account_id"`
	Amount int32 `json:"amount"`
}

type TransferTxResults struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

var transactionKey = struct{}{}


// 
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error) {
	var result TransferTxResults

	err := store.execTX(ctx, func(q *Queries) error {

		var err error

		// txName := ctx.Value(transactionKey)

		// fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx,  CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID: arg.ToAccountId,
			Amount: arg.Amount,
		})
		if err != nil {
			return  err
		}
	
		// fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount: -arg.Amount,
		})
		if err != nil {
			return  err
		}

		// fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount: arg.Amount,
		})
		if err != nil {
			return  err
		}


		if arg.FromAccountId < arg.ToAccountId {

			// fmt.Println(txName, "update account 1")
		result.FromAccount, err = q.AddAccountBalance(ctx , AddAccountBalanceParams{
			ID: arg.FromAccountId,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "update account 2")
		result.ToAccount, err = q.AddAccountBalance(ctx , AddAccountBalanceParams{
			ID: arg.ToAccountId,
			Amount: arg.Amount,
			
		})
		if err != nil {
			return err
		}

		} else {
			
		// fmt.Println(txName, "update account 2")
		result.ToAccount, err = q.AddAccountBalance(ctx , AddAccountBalanceParams{
			ID: arg.ToAccountId,
			Amount: arg.Amount,
			
		})
		if err != nil {
			return err
		}

		result.FromAccount, err = q.AddAccountBalance(ctx , AddAccountBalanceParams{
			ID: arg.FromAccountId,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}
		}
		


		return nil
	})

	return result, err

	
}