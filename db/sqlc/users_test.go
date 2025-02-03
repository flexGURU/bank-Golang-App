package db

import (
	"testing"

	"github.com/flexGURU/simplebank/utils"
	"github.com/stretchr/testify/require"
)




func CreateUserTestAccount(t *testing.T) User {

	arg := CreateUserParams {
		Username: utils.RandomOwner(),
		HashedPassword: "secret",
		FullName: utils.RandomOwner(),
		Email: utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)

	require.NotZero(t, user.CreatedAt)

	return user


	
}


func TestUserCreateAccount(t *testing.T)  {

	CreateUserTestAccount(t)
	
}

func TestGetUser(t *testing.T)  {

	user_account := CreateUserTestAccount(t)

	user, err := testQueries.GetUser(ctx, user_account.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user_account.Username, user.Username)
	require.Equal(t, user_account.Email, user.Email)
	require.Equal(t, user_account.HashedPassword, user.HashedPassword)
	require.Equal(t, user_account.FullName, user.FullName)

	require.NotZero(t, user.CreatedAt)

}
