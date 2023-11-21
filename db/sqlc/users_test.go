package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/silaselisha/bankapi/db/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateUsers(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	resp, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.Equal(t, user.Email, resp.Email)
	require.Equal(t, user.Fullname, resp.Fullname)
	require.Equal(t, user.Username, resp.Username)
	require.Equal(t, user.Password, resp.Password)

	require.WithinDuration(t, user.CreatedAt, resp.CreatedAt, 1*time.Second)
}

func createRandomUser(t *testing.T) *User {
	username, _ := utils.RandomString(6)
	firstName, _ := utils.RandomString(6)
	lastName, _ := utils.RandomString(6)

	args := CreateUserParams{
		Username: username,
		Fullname: fmt.Sprintf("%s %s", firstName, lastName),
		Email:    fmt.Sprintf("%s%d@gmail.com", username, utils.RandomAmount(1, 100)),
		Password: "secret",
	}

	user, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username: args.Username,
		Fullname: args.Fullname,
		Email:    args.Email,
		Password: args.Password,
	})

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.Fullname, user.Fullname)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.Password, user.Password)

	require.WithinDuration(t, time.Now(), user.CreatedAt, 1*time.Second)
	return &user
}
