-- name: CreateTransfer :one

INSERT INTO
    transfers (
        from_account_id,
        to_account_id,
        amount
    )
VALUES ($1, $2, $3) RETURNING *;

-- name: ListTransfers :many

SELECT *
FROM transfers
WHERE
    from_account_id = $1
    OR to_account_id = $2
ORDER BY id
OFFSET $3
LIMIT $4;

-- name: GetTransfer :one

SELECT * FROM transfers WHERE id = $1 LIMIT 1;