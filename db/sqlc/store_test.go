package db

import (
	"context"
	"database/sql"
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

	errs := make(chan error, 10)
	results := make(chan TransferTxResultsParams, 10)

	for i := 0; i < iterations; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				ToAccountId:   account2.ID,
				FromAccountId: account1.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// since it's 5 transaction running then iterations are gonna be 5
	for i := 0; i < iterations; i++ {
		result := <- results
		require.NotEmpty(t, result)
		err := <- errs
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
		require.Equal(t, sql.NullInt32{Int32: -amount, Valid: true}, fromEntry.Amount)
		require.WithinDuration(t, time.Now(), fromEntry.CreatedAt, time.Second * 1)
		
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, sql.NullInt32{Int32: amount, Valid: true}, toEntry.Amount)
		require.WithinDuration(t, time.Now(), toEntry.CreatedAt, time.Second * 1)
	}
}