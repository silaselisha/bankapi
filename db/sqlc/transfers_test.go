package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/silaselisha/bank-api/db/utils"
	"github.com/stretchr/testify/require"
)

func createSingleTransfer(t *testing.T, from, to Account) Transfer {

	args := CreateTransferParams{
		FromAccountID: sql.NullInt64{Int64: from.ID, Valid: true},
		ToAccountID: sql.NullInt64{Int64: to.ID, Valid: true},
		Amount: utils.GenerateAmount(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, sql.NullInt64{Int64: from.ID, Valid: true}, transfer.FromAccountID)
	require.Equal(t, sql.NullInt64{Int64: to.ID, Valid: true}, transfer.ToAccountID)


	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createSingleAccount(t)
	account2 := createSingleAccount(t)

	createSingleTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createSingleAccount(t)
	account2 := createSingleAccount(t)
	
	transfer := createSingleTransfer(t, account1, account2)

	transfer1, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer1)

	require.Equal(t, transfer.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer1.ToAccountID)
	require.Equal(t, transfer.Amount, transfer1.Amount)
	require.Equal(t, transfer.ID, transfer1.ID)

	require.WithinDuration(t, transfer.CreatedAt, transfer1.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	account1 := createSingleAccount(t)
	account2 := createSingleAccount(t)

	n := 6
	for i := 0; i < n; i++ {
		createSingleTransfer(t, account1, account2)
	}

	args := ListTransfersParams{
		FromAccountID: sql.NullInt64{Int64: account1.ID, Valid: true},
		ToAccountID: sql.NullInt64{Int64: account2.ID, Valid: true},
		Limit: 3,
		Offset: 3,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, sql.NullInt64{Int64: account1.ID, Valid: true}, transfer.FromAccountID)
		require.Equal(t, sql.NullInt64{Int64: account2.ID, Valid: true}, transfer.ToAccountID)
	}
}