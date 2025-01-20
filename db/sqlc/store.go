package db

import (
	"context"
	"database/sql"
	"fmt"
)

// the store will provide all functions to execute the database queries and transactions
type Store struct {
	db *sql.DB
	*Queries

}

// A constructor that create a new instance of the store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

// execTx exuctes a function within a database transaction
func (store *Store) execTX(ctx context.Context, fn func(*Queries) error) error {

	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	
	
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); err != nil {
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


// 
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error) {
	var result TransferTxResults

	err := store.execTX(ctx, func(q *Queries) error {

		var err error
		result.Transfer, err = q.CreateTransfer(ctx,  CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID: arg.ToAccountId,
			Amount: float64(arg.Amount),
		})
		if err != nil {
			return  err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount: -arg.Amount,
		})
		if err != nil {
			return  err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount: arg.Amount,
		})
		if err != nil {
			return  err
		}


		return nil
	})

	return result, err

	
}