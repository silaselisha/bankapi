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

var txKey = struct{}{}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// ** begin a transaction
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q) 

	if err != nil {
		if rollback_err := tx.Rollback(); rollback_err != nil {
			return fmt.Errorf("transaction error: %v | rollback error: %v", err, rollback_err)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxResult struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to-account"`
	FromEntry   Entry    `json:"from-entry"`
	ToEntry     Entry    `json:"to_entry"`
	Transfer    Transfer `json:"transfer"`
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	txValue := ctx.Value(txKey)
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		fmt.Printf("Tx>> %v create transfer\n", txValue)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: sql.NullInt64{Int64: args.FromAccountID, Valid: true},
			ToAccountID:   sql.NullInt64{Int64: args.ToAccountID, Valid: true},
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Tx>> %v create from entry\n", txValue)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: sql.NullInt64{Int64: args.FromAccountID, Valid: true},
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Tx>> %v create to entry\n", txValue)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: sql.NullInt64{Int64: args.ToAccountID, Valid: true},
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		
		if args.FromAccountID < args.ToAccountID {
			account1, err := q.GetAccountForUpdate(ctx, args.FromAccountID)
	
			if err != nil {
				return err
			}

			fmt.Printf("Tx>> %v update account\n", txValue)
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      account1.ID,
				Balance: account1.Balance - args.Amount,
			})
			if err != nil {
				return err
			}
	
			account2, err := q.GetAccountForUpdate(ctx, args.ToAccountID)
	
			if err != nil {
				return err
			}
	
			fmt.Printf("Tx>> %v update account\n", txValue)
			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      account2.ID,
				Balance: account2.Balance + args.Amount,
			})
			if err != nil {
				return err
			}
		}else {
			account2, err := q.GetAccountForUpdate(ctx, args.ToAccountID)
	
			if err != nil {
				return err
			}
	
			fmt.Printf("Tx>> %v update account\n", txValue)
			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      account2.ID,
				Balance: account2.Balance + args.Amount,
			})
			if err != nil {
				return err
			}
			
			account1, err := q.GetAccountForUpdate(ctx, args.FromAccountID)
	
			if err != nil {
				return err
			}

			fmt.Printf("Tx>> %v update account\n", txValue)
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      account1.ID,
				Balance: account1.Balance - args.Amount,
			})
			if err != nil {
				return err
			}

		}

		return nil
	})
	return result, err
}