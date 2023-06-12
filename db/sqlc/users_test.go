package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/silaselisha/bank-api/db/utils"
	"github.com/stretchr/testify/require"
)

func createSingleUser(t *testing.T) User {
	firstName := utils.GenerateFirstName()

	args := CreateUserParams{
		FirstName: firstName,
		LastName:  utils.GenerateLastName(),
		Gender:    utils.GenerateGender(),
		Email:     fmt.Sprintf("%s@gmail.com", firstName),
		Password:  "",
	}

	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.LastName, user.LastName)
	require.Equal(t, args.FirstName, user.FirstName)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCerateUser(t *testing.T) {
	createSingleUser(t)
}

func TestGetUser(t *testing.T) {
	testUser := createSingleUser(t)

	user, err := testQueries.GetUser(context.Background(), testUser.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, testUser.ID, user.ID)
	require.Equal(t, testUser.Email, user.Email)
	require.Equal(t, testUser.LastName, user.LastName)
	require.Equal(t, testUser.FirstName, user.FirstName)
	require.WithinDuration(t, testUser.CreatedAt, user.CreatedAt, time.Second)
}

func TestListUser(t *testing.T) {
	n := 6

	for i := 0; i < n; i++ {
		createSingleUser(t)
	}

	args := ListUsersParams {
		Limit: 4,
		Offset: 2,
	}
	users, err := testQueries.ListUsers(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	
	for _, user := range users {
		require.NotEmpty(t, user)
		require.NotZero(t, user.ID)
		require.NotZero(t, user.CreatedAt)
	}
}