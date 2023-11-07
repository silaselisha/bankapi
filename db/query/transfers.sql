-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
) VALUES ($1, $2, $3) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers WHERE from_account_id = $1 AND to_account_id = $2 OFFSET $3 LIMIT $4;