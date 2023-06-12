package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/silaselisha/bank-api/db/utils"
	"github.com/stretchr/testify/require"
)

func createSingleAccount(t *testing.T) Account {
	user := createSingleUser(t)
	args := CreateAccountParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Gender:    utils.GenerateGender(),
		Balance:   utils.GenerateAmount(),
		Currency:  utils.GenerateCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.FirstName, args.FirstName)
	require.Equal(t, account.LastName, args.LastName)
	require.Equal(t, account.Gender, args.Gender)
	require.Equal(t, account.Balance, args.Balance)
	require.Equal(t, account.Currency, args.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createSingleAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createSingleAccount(t)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.FirstName, account1.FirstName)
	require.Equal(t, account2.LastName, account1.LastName)
	require.Equal(t, account2.Gender, account1.Gender)
	require.Equal(t, account2.Balance, account1.Balance)

	require.WithinDuration(t, account2.CreatedAt, account1.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createSingleAccount(t)

	args := UpdateAccountParams{
		ID: account1.ID,
		Balance: utils.GenerateAmount(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account2.Balance, args.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account := createSingleAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

	account1, err := testQueries.GetAccount(context.Background(), account.ID)

	require.Error(t, err)
	require.Empty(t, account1)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 5; i++ {
		createSingleAccount(t)
	}

	args := ListAccountsParams{
		Limit: 3,
		Offset: 3,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	
	k := len(accounts)
	require.Equal(t, k, int(args.Limit))

	for _, account := range accounts {
		require.NotEmpty(t, account)
		fmt.Println(account)
	}
}