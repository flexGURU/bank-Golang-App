package token

import (
	"testing"
	"time"

	"github.com/flexGURU/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))

	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)


	JWTpayload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, JWTpayload)

	require.NotZero(t, JWTpayload.ID)
	require.Equal(t, username, JWTpayload.Username)
	require.WithinDuration(t, issuedAt, JWTpayload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, JWTpayload.ExpiredAt, time.Second)
	
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))

	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute
	token, payload, err := maker.CreateToken(username, -duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)


	payload, err = maker.VerifyToken(token)

	require.Error(t, err)
	// require.EqualError(t, err, fmt.Sprintf("expired token"))
	require.Nil(t, payload)
}