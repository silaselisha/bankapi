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
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/silaselisha/bankapi/db/mock"
	db "github.com/silaselisha/bankapi/db/sqlc"
	"github.com/silaselisha/bankapi/db/utils"
	"github.com/silaselisha/bankapi/token"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	testCases := []struct {
		name    string
		id      int64
		stub    func(store *mockdb.MockStore)
		setAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		check   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			id:   account.ID,
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationType, account.Owner, 30*time.Minute)
			},
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
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationType, account.Owner, 30*time.Minute)
			},
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
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationType, account.Owner, 30*time.Minute)
			},
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
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationType, account.Owner, 2*time.Minute)
			},
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
			key := "53d36d97c0dbbec983acdf2ae98746d7"
			tokenMaker, err := token.NewJwtMaker(key)

			tsc.setAuth(t, request, tokenMaker)
			fmt.Println(request.Header)
			server.router.ServeHTTP(recoder, request)
			require.NoError(t, err)
			tsc.check(t, recoder)
		})
	}
}

func createRandomAccount(t *testing.T) db.Account {
	user, _ := createRandomUser(t)
	return db.Account{
		ID:       int64(utils.RandomAmount(1, 100)),
		Owner:    user.Username,
		Balance:  int32(utils.RandomAmount(100, 1000)),
		Currency: utils.RandomCurrency(),
	}
}

func createRandomUser(t *testing.T) (db.User, string) {
	username, _ := utils.RandomString(6)
	firstName, _ := utils.RandomString(6)
	lastName, _ := utils.RandomString(6)
	password, _ := utils.RandomString(8)
	hashedPassword, err := utils.GenerateHashedPassword(password)

	require.NoError(t, err)
	return db.User{
		Username:  username,
		Fullname:  fmt.Sprintf("%s %s", firstName, lastName),
		Email:     fmt.Sprintf("%s%d@gmail.com", username, utils.RandomAmount(1, 100)),
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}, password
}

func checkBodyResponse(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var res db.Account
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	require.Equal(t, account, res)
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(AuthorizationHeaderKey, authorizationHeader)
}
