package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/flexGURU/simplebank/db/mock"
	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {

    account := randomAccount()

    testCases := []struct{
        name string
        accountID int32
        stub func(store *mockdb.MockStore)
        checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
    }{
        {
            name: "status OK",
            accountID: account.ID,
            stub: func(store *mockdb.MockStore){
                store.EXPECT().
                    GetAccount(gomock.Any(), gomock.Eq(account.ID)).
                    Times(1).   
                    Return(account, nil)

            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
                require.Equal(t,  http.StatusOK, recorder.Code)
                checkResponseBody(t, recorder.Body, account)

            },
        },
        {
            name: "Not Found",
            accountID: account.ID,
            stub: func(store *mockdb.MockStore){
                store.EXPECT().
                GetAccount(gomock.Any(), gomock.Eq(account.ID)).
                Times(1).
                Return(db.Account{}, sql.ErrNoRows)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
                require.Equal(t,  http.StatusNotFound, recorder.Code)

            },
        },
        {
            name: "InternalError",
            accountID: account.ID,
            stub: func(store *mockdb.MockStore){
                store.EXPECT().
                GetAccount(gomock.Any(), gomock.Eq(account.ID)).
                Times(1).
                Return(db.Account{}, sql.ErrConnDone)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
                require.Equal(t,  http.StatusInternalServerError, recorder.Code)

            },
        },
        {
            name: "InvalidID",
            accountID: 0,
            stub: func(store *mockdb.MockStore){
                store.EXPECT().
                GetAccount(gomock.Any(), gomock.Any()).
                Times(0)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
                require.Equal(t,  http.StatusBadRequest, recorder.Code)

            },
        },



    
    }

    for i := range testCases {
        tc := testCases[i]
        

        t.Run(tc.name, func(t *testing.T){
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
        
        
            store := mockdb.NewMockStore(ctrl)
            tc.stub(store)
        
        
            // creating a server for test and send the request
            server := newTestServer(t, store)
            require.NotEmpty(t, server)
            recorder := httptest.NewRecorder()
        
            url := fmt.Sprintf("/getaccount/%d", tc.accountID)
        
            request, err := http.NewRequest(http.MethodPost, url, nil)
            require.NoError(t, err)
        
        
            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(t, recorder)
        })
		


        
    }

}


func randomAccount() db.Account {
    return db.Account{
        ID: utils.RandomInt(1,100),
        Owner: utils.RandomOwner(),
        Balance: utils.RandomMoney(),
        Currency: utils.RandomCurrency(),
    }


}


func checkResponseBody(t *testing.T, body *bytes.Buffer, account db.Account)  {

    data, err := io.ReadAll(body)
    require.NoError(t, err)
    
    var gotAccount db.Account

    err = json.Unmarshal(data, &gotAccount)
    require.NoError(t, err)
    require.Equal(t, account, gotAccount)
    
}