package api

import (
	"os"
	"testing"
	"time"

	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymmetricKey: utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, nil)
	require.NoError(t, err)
	require.NotEmpty(t, server)

	return server
}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())

}