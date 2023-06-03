package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(conn)

	account1 := createSingleAccount(t)
	account2 := createSingleAccount(t)

	amount := int64(25)
	n := 4
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx%v", i+1)
		ctx := context.WithValue(context.Background(), txKey, txName)
		go func() {
			transfer, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- transfer
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		transfer := result.Transfer

		require.NotEmpty(t, transfer)

		require.Equal(t, sql.NullInt64{Int64: account1.ID, Valid: true}, transfer.FromAccountID)
		require.Equal(t, sql.NullInt64{Int64: account2.ID, Valid: true}, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromentry := result.FromEntry
		require.NotEmpty(t, fromentry)

		require.Equal(t, sql.NullInt64{Int64: account1.ID, Valid: true}, fromentry.AccountID)
		require.Equal(t, -amount, fromentry.Amount)

		require.NotZero(t, fromentry.ID)
		require.NotZero(t, fromentry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromentry.ID)
		require.NoError(t, err)

		toentry := result.ToEntry
		require.NotEmpty(t, toentry)

		require.Equal(t, sql.NullInt64{Int64: account2.ID, Valid: true}, toentry.AccountID)
		require.Equal(t, amount, toentry.Amount)

		require.NotZero(t, toentry.ID)
		require.NotZero(t, toentry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toentry.ID)
		require.NoError(t, err)

		// TODO: [GBA-2] Test user's account balance
		fromaccount := result.FromAccount

		require.NotEmpty(t, fromaccount)
		require.Equal(t, sql.NullInt64{Int64: account1.ID, Valid: true}, sql.NullInt64{Int64: fromaccount.ID, Valid: true})

		toaccount := result.ToAccount

		require.NotEmpty(t, toaccount)
		require.Equal(t, sql.NullInt64{Int64: account2.ID, Valid: true}, sql.NullInt64{Int64: toaccount.ID, Valid: true})

		diff1 := account1.Balance - fromaccount.Balance
		diff2 := toaccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	require.Equal(t, account1.Balance - int64(n) * amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n) * amount, updatedAccount2.Balance)
}

func TestTransferTxAcc(t *testing.T){
	store := NewStore(conn)
	account1 := createSingleAccount(t)
	account2 := createSingleAccount(t)

	amount := int64(25)
	n := 6

	errs := make(chan error)

	/*
	*@TODO: Handle transaction that are bi-directional
	*/
	fmt.Printf("Before Account one>> %d\n", account1.Balance)
	fmt.Printf("Before Account two>> %d\n", account2.Balance)
	for i := 0; i < n; i++ {
		FromAccountID := account1.ID
		ToAccountID := account2.ID
		
		txName := fmt.Sprintf("tx%d", i+1)
		ctx := context.WithValue(context.Background(), txKey, txName)

		if i % 2 == 1 {
			FromAccountID = account2.ID
			ToAccountID = account1.ID
		}

		go func () {
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: FromAccountID,
				ToAccountID: ToAccountID,
				Amount: amount,
			})

			errs <- err
		} ()
	}	

	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
	}
	// **TODO: remove all prints

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
	fmt.Printf("After account one>>  %d\n", updatedAccount1.Balance)
	fmt.Printf("After account one>>  %d\n", updatedAccount2.Balance)
}