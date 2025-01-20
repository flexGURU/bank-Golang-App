package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(dbConn)

	account1 := createTestAccount(t)
	account2 := createTestAccount(t)


	// run database transactions using concurrency
	n := 5
	amount := 60

	results := make(chan TransferTxResults)
	errs := make(chan error)


	for i := 0; i < n; i ++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId: account2.ID,
				Amount: int32(amount),
			})

			errs <- err 
			results <- result
			
		}()
	}


	for i := 0; i < n; i ++ {

		err := <- errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, amount, result.Transfer.Amount)

		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)



	}




}