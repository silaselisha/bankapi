package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/silaselisha/bankapi/db/mock"
	db "github.com/silaselisha/bankapi/db/sqlc"
	"github.com/silaselisha/bankapi/db/utils"
	"github.com/silaselisha/bankapi/token"
	"github.com/stretchr/testify/require"
)

func TestTransfers(t *testing.T) {
	account1 := createRandomAcc()
	account2 := createRandomAcc()
	amount := int32(20)

	testCases := []struct {
		name  string
		body  gin.H
		setAuth func(t *testing.T, request *http.Request, maker token.Maker)
		stub  func(store *mockdb.MockStore)
		check func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "200 ok",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.USD,
			},
			setAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, AuthorizationType, account1.Owner, 30 * time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)

				args := db.TransferTxParams{
					FromAccountId: account1.ID,
					ToAccountId:   account2.ID,
					Amount:        amount,
				}

				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(args)).Times(1)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "404 bad request from account",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.EUR,
			},
			setAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, AuthorizationType, account1.Owner, 30 * time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
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
			tsc.stub(store)

			server := NewServer(store)
			url := "/transfers"
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tsc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			key := "53d36d97c0dbbec983acdf2ae98746d7"
			maker, err := token.NewJwtMaker(key)
			require.NoError(t, err)
			tsc.setAuth(t, request, maker)
			server.router.ServeHTTP(recorder, request)
		})
	}
}

func createRandomAcc() db.Account {
	owner, _ := utils.RandomString(6)

	return db.Account{
		ID:        int64(utils.RandomAmount(1, 100)),
		Owner:     owner,
		Balance:   int32(utils.RandomAmount(100, 2000)),
		Currency:  "USD",
		CreatedAt: time.Now(),
	}
}
