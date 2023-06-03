package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/silaselisha/bank-api/db/utils"
	"github.com/stretchr/testify/require"
)

func createSingleEntry(t *testing.T, account Account) Entry {

	args := CreateEntryParams{
		AccountID: sql.NullInt64{Int64: account.ID, Valid: true},
		Amount: utils.GenerateAmount(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, sql.NullInt64{Int64: account.ID, Valid: true}, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	return entry
}
func TestCreateEntry(t *testing.T) {
	account := createSingleAccount(t)
	createSingleEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createSingleAccount(t)
	entry := createSingleEntry(t, account)

	entry1, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry1)
	
	require.Equal(t, entry.AccountID, entry1.AccountID)
	require.Equal(t, entry.Amount, entry1.Amount)

	require.NotZero(t, entry1.ID)
	require.NotZero(t, entry1.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account := createSingleAccount(t)
	n := 6

	for i := 0; i < n; i++ {
		createSingleEntry(t, account)
	}

	args := ListEntriesParams{
		AccountID: sql.NullInt64{Int64: account.ID, Valid: true},
		Limit: 3,
		Offset: 3,
	}


	entries, err := testQueries.ListEntries(context.Background(), args)
	fmt.Println(entries)

	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, sql.NullInt64{Int64: account.ID, Valid: true}, entry.AccountID)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)
	}
}