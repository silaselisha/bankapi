package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/silaselisha/bankapi/database/utils"
	"github.com/stretchr/testify/require"
)

func TestAccounts(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccountById(t *testing.T) {
	test_account := createRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), test_account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, test_account.ID, account.ID)
	require.Equal(t, test_account.Owner, account.Owner)
	require.Equal(t, test_account.Balance, account.Balance)
	require.Equal(t, test_account.Currency, account.Currency)
	require.WithinDuration(t, test_account.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccountBalance(t *testing.T) {
	account := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      account.ID,
		Amount: int32(utils.RandomAmount(10, 500)),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, account.Owner, updatedAccount.Owner)
	require.Equal(t, args.Amount, updatedAccount.Balance)
	require.WithinDuration(t, account.CreatedAt, updatedAccount.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomAccount(t)
	}

	agrs := ListAccountsParams{
		Offset: 1,
		Limit:  5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), agrs)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.True(t, len(accounts) <= 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestDeleteAccountById(t *testing.T) {
	test_account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), test_account.ID)
	require.NoError(t, err)
	account, err := testQueries.GetAccount(context.Background(), test_account.ID)

	require.Equal(t, Account{ID: 0, Owner: "", Balance: 0, Currency: "", CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, account)
	require.Equal(t, sql.ErrNoRows, err)
}

func createRandomAccount(t *testing.T) *Account {
	owner, err := utils.RandomString(6)
	require.NoError(t, err)

	args := CreateAccountParams{
		Owner:    owner,
		Balance:  int32(utils.RandomAmount(100, 2000000)),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Currency, account.Currency)
	require.Equal(t, args.Balance, account.Balance)
	require.WithinDuration(t, time.Now(), account.CreatedAt, time.Second*2)
	return &account
}
