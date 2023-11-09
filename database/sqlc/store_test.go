package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(conn)
	const iterations int = 5
	const amount int32 = 100

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("account 1 balance: ", account1.Balance)
	fmt.Println("account 2 balance: ", account2.Balance)
	exists := make(map[int]bool)
	errs := make(chan error, 10)
	results := make(chan TransferTxResultsParams, 10)

	for i := 0; i < iterations; i++ {
		value := fmt.Sprintf("tx::%v", i+1)
		ctx := context.WithValue(context.Background(), txKey, value)
		go func() {
			result, err := store.TransferTx(ctx, TransferTxParams{
				ToAccountId:   account2.ID,
				FromAccountId: account1.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < iterations; i++ {
		result := <-results
		require.NotEmpty(t, result)
		err := <-errs
		require.NoError(t, err)

		// working on each transaction data
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.WithinDuration(t, time.Now(), transfer.CreatedAt, time.Second)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.ID)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.WithinDuration(t, time.Now(), fromEntry.CreatedAt, time.Second*1)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.WithinDuration(t, time.Now(), toEntry.CreatedAt, time.Second*1)
		require.Equal(t, amount, toEntry.Amount)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// account balance updates
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		fmt.Println("from account balance: ", fromAccount.Balance)
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		fmt.Println("to account balance: ", toAccount.Balance)

		// -> A1 = 200 | A2 160
		// -> A1=200 - 20 | A2= 160 + 20
		// -> diff1 = 200 - 180 -> 20
		// -> diff2 = 180 - 160 -> 20
		// -> 20 / 20 = 1
		// -> diff1 % iterations = 0 | 20 % 5 = 0
		diff2 := toAccount.Balance - account2.Balance
		diff1 := account1.Balance - fromAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1/amount > 0)
		require.True(t, diff1%int32(iterations) == 0)
		key := int(diff1 / amount)
		require.NotContains(t, exists, key)
		exists[key] = true

		require.Equal(t, account1.Balance-(int32(iterations)*amount), fromAccount.Balance)
		require.Equal(t, account2.Balance+(int32(iterations)*amount), toAccount.Balance)
	}

	fromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance-(int32(iterations)*amount), fromAccount.Balance)

	toAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Balance+(int32(iterations)*amount), toAccount.Balance)
}
func TestTransferAlternateTx(t *testing.T) {
	store := NewStore(conn)
	const iterations int = 4
	const amount int32 = 100

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	errs := make(chan error, 10)
	results := make(chan TransferTxResultsParams, 10)

	for i := 0; i < iterations; i++ {
		toAccountId := account2.ID
		fromAccountId := account1.ID

		if i % 2 == 1 {
			toAccountId = account1.ID
			fromAccountId = account2.ID
		}

		value := fmt.Sprintf("tx::%v", i+1)
		ctx := context.WithValue(context.Background(), txKey, value)
		go func() {
			result, err := store.TransferTx(ctx, TransferTxParams{
				ToAccountId:   toAccountId,
				FromAccountId: fromAccountId,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < iterations; i++ {
		result := <-results
		require.NotEmpty(t, result)
		err := <-errs
		require.NoError(t, err)
	}

	fromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	toAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, (account1.Balance - fromAccount.Balance), (toAccount.Balance - account2.Balance))
}