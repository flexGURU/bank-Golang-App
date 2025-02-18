package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/flexGURU/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuth(
	t *testing.T,
	tokenMaker token.Maker,
	username string,
	request *http.Request,
	duration time.Duration,
	authorizationTypeBearer string,
)  {

	token, payload,  err := tokenMaker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
	request.Header.Set(authorizationHeaderKey, authHeader)
	
}

func TestAuthMiddleware(t *testing.T) {

	testCases := []struct{
		name string
		setUpAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// happy test case
		{
			name: "OK",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, tokenMaker, "user", request, time.Minute, authorizationTypeBearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		//  No auth
		{
			name: "No authorisation",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		//  unsupported auth
		{
			name: "Unsupported authorisation",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, tokenMaker, "user", request, time.Minute, "unsupported")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		//  Invalid auth format
		{
			name: "Unsupported authorisation",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, tokenMaker, "user", request, time.Minute, "")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		//  Iexpired access token
		{
			name: "expired token",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, tokenMaker, "user", request, -time.Minute, authorizationTypeBearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		
	}


	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				} ,
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setUpAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t,recorder)
		})
	}


}