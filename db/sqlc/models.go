// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"time"
)

type Account struct {
	ID        int64     `db:"id"`
	Owner     string    `db:"owner"`
	Balance   int32     `db:"balance"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
}

type Entry struct {
	ID        int64     `db:"id"`
	AccountID int64     `db:"account_id"`
	Amount    int32     `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

type Transfer struct {
	ID            int64     `db:"id"`
	FromAccountID int64     `db:"from_account_id"`
	ToAccountID   int64     `db:"to_account_id"`
	Amount        int32     `db:"amount"`
	CreatedAt     time.Time `db:"created_at"`
}

type User struct {
	Username  string    `db:"username"`
	Fullname  string    `db:"fullname"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
