package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestTransferTx(t *testing.T) {

	store := NewStore(dbConn)

	// Creatng test accounts
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)
	fmt.Println("account balances before transaction", account1.Balance, account2.Balance)


	// run database transactions using concurrency
	n := 10
	amount := int32(10)

	// results := make(chan TransferTxResults)
	errs := make(chan error)

	for i := 0; i < n; i ++ {

		fromAccountID := account1.ID
		toAccountID := account2.ID

		if (i % 2 == 1) {		
			fromAccountID = account2.ID
			toAccountID = account1.ID

		}
		
		transactionName := fmt.Sprintf("transaction %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), transactionKey, transactionName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: fromAccountID,
				ToAccountId: toAccountID,
				Amount: amount,
			})

			errs <- err 
			// results <- result
			
		}()
	}


	// existed := make(map[int]bool)
	for i := 0; i < n; i ++ {

		err := <- errs
		require.NoError(t, err)

		// result := <- results
		// require.NotEmpty(t, result)

		// require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		// require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		// require.Equal(t, amount, result.Transfer.Amount)

		// require.NotZero(t, result.Transfer.CreatedAt)

		// _, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		// require.NoError(t, err)

		// fromEntry := result.FromEntry

		// require.NotEmpty(t, fromEntry)
		
		// require.Equal(t, account1.ID, fromEntry.AccountID)
		// require.Equal(t, -amount, fromEntry.Amount)

		// require.NotZero(t, fromEntry.ID)
		// require.NotZero(t, fromEntry.CreatedAt)

		// _, err = store.GetEntry(context.Background(), fromEntry.ID)

		// require.NoError(t, err)


		// toEntry := result.ToEntry

		// require.NotEmpty(t, toEntry)
		
		// require.Equal(t, account2.ID, toEntry.AccountID)
		// require.Equal(t, amount, toEntry.Amount)

		// require.NotZero(t, toEntry.ID)
		// require.NotZero(t, toEntry.CreatedAt)

		// _, err = store.GetEntry(context.Background(), toEntry.ID)

		// require.NoError(t, err)


		// // Check accounts
		// fromAccount := result.FromAccount
		// require.NotEmpty(t, fromAccount)
		// require.Equal(t, account1.ID, fromAccount.ID)


		// toAccount := result.ToAccount
		// require.NotEmpty(t, toAccount)
		// require.Equal(t, account2.ID, toAccount.ID)

		// fmt.Println("account balances after each transaction ",fromAccount.Balance, toAccount.Balance)

		// // checking acount balances
		// diff1 := account1.Balance - fromAccount.Balance
		// diff2 := toAccount.Balance - account2.Balance
		// require.Equal(t, diff1, diff2)
		// require.True(t,diff1 > 0)
		// require.True(t, diff1 % amount == 0)


		// k := int(diff1/amount)
		// require.True(t, k>=1 && k <= n)

		// require.NotContains(t, existed, k)
		// existed[k] = true


		// checking tthe final updated accounts
		updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)

		require.NoError(t, err)
		require.NotEmpty(t, updatedAccount1)


		// checking tthe final updated accounts
		updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)

		require.NoError(t, err)
		require.NotEmpty(t, updatedAccount2)

		fmt.Println("account balances after transaction", updatedAccount1.Balance, updatedAccount2.Balance)

		require.Equal(t, account1.Balance, updatedAccount1.Balance )
		require.Equal(t, account2.Balance, updatedAccount2.Balance )

	}




}