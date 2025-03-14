package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/flexGURU/simplebank/utils"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()


func createTestAccount(t *testing.T) Account {
	user := CreateUserTestAccount(t)

	arg := CreateAccountParams{
		Owner: user.Username,
		Balance: (utils.RandomMoney()),
		Currency: utils.RandomCurrency(),
	}


	account, err := testQueries.CreateAccount(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account

}

func TestCreateAccount(t *testing.T) {

	createTestAccount(t)


}

func TestGetAccount(t *testing.T)  {

	account1 := createTestAccount(t)

	account2, err := testQueries.GetAccount(ctx, account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)


	
}

func TestUpdateAccount(t *testing.T) {
	account1 := createTestAccount(t)

	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(ctx, arg )

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	fmt.Println(account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account2.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Currency, account2.Currency)



	
}