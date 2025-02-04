package auth

import (
	"testing"

	"github.com/flexGURU/simplebank/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {

	password := utils.RandomString(6)

	hashedpwd, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedpwd)

	err = ComparePassword(hashedpwd, password)

	require.NoError(t, err)

	// testing for a wrong password
	wrongPassword := utils.RandomString(6)

	err = ComparePassword(hashedpwd, wrongPassword)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())




}