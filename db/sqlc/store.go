package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResultsParams, error)
}
type SQLstore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLstore{
		Queries: New(db),
		db:      db,
	}
}

var txKey struct{} = struct{}{}

func (store *SQLstore) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (store *SQLstore) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResultsParams, error) {
	var results TransferTxResultsParams
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		value := ctx.Value(txKey)
		fmt.Printf("create transfer: %v\n", value)
		results.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountId,
			ToAccountID:   args.ToAccountId,
			Amount:        args.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Printf("create entry 1: %v\n", value)
		results.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountId,
			Amount:    -args.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Printf("create entry 2: %v\n", value)
		results.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountId,
			Amount:    args.Amount,
		})

		if err != nil {
			return err
		}

		if args.FromAccountId < args.ToAccountId {
			results.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:     args.FromAccountId,
				Amount: -args.Amount,
			})
			if err != nil {
				return err
			}

			fmt.Printf("update receiver balance: %v\n", value)
			results.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:     args.ToAccountId,
				Amount: +args.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("update receiver balance: %v\n", value)
			results.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:     args.ToAccountId,
				Amount: +args.Amount,
			})
			if err != nil {
				return err
			}

			fmt.Printf("update senders balance: %v\n", value)
			results.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:     args.FromAccountId,
				Amount: -args.Amount,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return results, err
}
