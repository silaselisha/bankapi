package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		if rlberr := tx.Rollback(); rlberr != nil {
			fmt.Println("rolback error...")
			return rlberr
		}
		fmt.Println("begin transaction error...")
		return err
	}

	queries := New(tx)
	err = fn(queries)
	if err != nil {
		fmt.Println("transaction error...")
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	ToAccountId   int64 `json:"to_acccount_id"`
	FromAccountId int64 `json:"from_account_id"`
	Amount        int32 `json:"amount"`
}

type TransferTxResultsParams struct {
	ToEntry     Entry    `json:"to_entry"`
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToAccount   Account  `json:"to_account"`
	FromAccount Account  `json:"from_account"`
}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResultsParams, error) {
	var results TransferTxResultsParams
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		results.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountId,
			ToAccountID:   args.ToAccountId,
			Amount:        args.Amount,
		})

		if err != nil {
			return err
		}

		results.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountId,
			Amount:    sql.NullInt32{Int32: -args.Amount, Valid: true},
		})

		if err != nil {
			return err
		}

		results.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountId,
			Amount: sql.NullInt32{Int32: args.Amount, Valid: true},
		})

		if err != nil {
			return err
		}

		return nil
	})

	return results, err
}
