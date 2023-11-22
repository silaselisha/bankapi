package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/silaselisha/bankapi/db/mock"
	db "github.com/silaselisha/bankapi/db/sqlc"
	"github.com/silaselisha/bankapi/db/utils"
	"github.com/stretchr/testify/require"
)

func TestUsers(t *testing.T) {
	user, password := createRandomUser(t)
	hashedPassword, err := utils.GenerateHashedPassword(password)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		username string
		body     gin.H
		stub     func(store *mockdb.MockStore)
		check    func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "Ok",
			username: user.Username,
			body: gin.H{
				"username": user.Username,
				"fullname": user.Fullname,
				"email":    user.Email,
				"password": password,
			},
			stub: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					Fullname: user.Fullname,
					Email:    user.Email,
					Password: hashedPassword,
				}

				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(args)).Times(1).Return(user, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tsc := testCases[i]
		t.Run(tsc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server := NewServer(store)

			data, err := json.Marshal(tsc.body)
			require.NoError(t, err)

			url := "/users"
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
		})
	}
}
