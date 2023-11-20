package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/silaselisha/bankapi/database/mock"
	database "github.com/silaselisha/bankapi/database/sqlc"
	"github.com/silaselisha/bankapi/database/utils"
	"github.com/stretchr/testify/require"
)

func TestTransfers(t *testing.T) {
	account1 := createRandomAcc()
	account2 := createRandomAcc()
	amount := int32(20)

	fmt.Println(account1, account2)

	testCases := []struct {
		name string
		body gin.H
		stub func(store *mockdb.MockStore)
	}{
		{
			name: "Ok",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":         "USD",
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)


				args := database.TransferTxParams{
					FromAccountId: account1.ID,
					ToAccountId:   account2.ID,
					Amount:        amount,
				}

				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(args)).Times(1)
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
			server.router.ServeHTTP(recorder, request)
		})
	}
}

func createRandomAcc() database.Account {
	owner, _ := utils.RandomString(6)

	return database.Account{
		ID:       int64(utils.RandomAmount(1, 100)),
		Owner:    owner,
		Balance:  int32(utils.RandomAmount(100, 2000)),
		Currency: "USD",
		CreatedAt: time.Now(),
	}
}
