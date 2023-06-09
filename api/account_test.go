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

	"github.com/golang/mock/gomock"
	mockdb "github.com/silaselisha/bank-api/db/mock"
	db "github.com/silaselisha/bank-api/db/sqlc"
	"github.com/silaselisha/bank-api/db/utils"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := createSingleAccount()

	testCases := []struct {
		name       string
		id         int64
		createStub func(store *mockdb.MockStore)
		compareAcc func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			id:   account.ID,
			createStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			compareAcc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				compareAccounts(t, recorder.Body, account)
			},
		},
		{
			name: "Ok",
			id:   account.ID,
			createStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			compareAcc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				compareAccounts(t, recorder.Body, account)
			},
		},
		{
			name: "Not Found",
			id:   account.ID,
			createStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			compareAcc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			id:   account.ID,
			createStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			compareAcc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Bad Get Request",
			id:   0,
			createStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			compareAcc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tsc := testCases[i]

		t.Run(tsc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
		
			store := mockdb.NewMockStore(ctrl)
			tsc.createStub(store)
		
			server := NewServer(store)
			recorder := httptest.NewRecorder()
		
			url := fmt.Sprintf("/accounts/%d", tsc.id)
		
			request := httptest.NewRequest(http.MethodGet, url, nil)
			server.router.ServeHTTP(recorder, request)
		
			tsc.compareAcc(t, recorder)
		})
	}
}

func createSingleAccount() db.Account {
	return db.Account{
		ID:        utils.RandomInteger(1, 20),
		FirstName: utils.GenerateFirstName(),
		LastName:  utils.GenerateLastName(),
		Gender:    utils.GenerateGender(),
		Balance:   utils.GenerateAmount(),
		Currency:  utils.GenerateCurrency(),
	}
}

func compareAccounts(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var parsedData db.Account
	err = json.Unmarshal(data, &parsedData)
	require.NoError(t, err)
	require.Equal(t, account, parsedData)
}
