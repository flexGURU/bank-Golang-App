package db

import (
	"context"
	"testing"

	"github.com/flexGURU/simplebank/utils"
	"github.com/stretchr/testify/require"
)


func TestCreateAccount(t *testing.T) {

	arg := CreateAccountParams{
		Owner: utils.RandomOwner(),
		Balance: float64(utils.RandomMoney()),
		Currency: utils.RandomCurrency(),
	}

	ctx := context.Background()

	account, err := testQueries.CreateAccount(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)









}