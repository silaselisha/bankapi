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
	mockdb "github.com/silaselisha/bankapi/db/mock"
	db "github.com/silaselisha/bankapi/db/sqlc"
	"github.com/silaselisha/bankapi/db/utils"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	testCases := []struct {
		name  string
		id    int64
		stub  func(store *mockdb.MockStore)
		check func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			id:   account.ID,
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkBodyResponse(t, recorder.Body, account)
			},
		},
		{
			name: "Not found",
			id:   account.ID,
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			id:   account.ID,
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid request",
			id:   0,
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tsc := range testCases {
		t.Run(tsc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockstore := mockdb.NewMockStore(ctrl)
			tsc.stub(mockstore)
			server := NewServer(mockstore)
			recoder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tsc.id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recoder, request)
			tsc.check(t, recoder)
		})
	}
}

func createRandomAccount(t *testing.T) db.Account {
	user, err := utils.RandomString(6)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return db.Account{
		ID:       int64(utils.RandomAmount(1, 100)),
		Owner:    user,
		Balance:  int32(utils.RandomAmount(100, 1000)),
		Currency: utils.RandomCurrency(),
	}
}

func checkBodyResponse(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var res db.Account
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	require.Equal(t, account, res)
}