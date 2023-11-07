-- name: CreateEntry :one

INSERT INTO
    entries (account_id, amount)
VALUES ($1, $2)
RETURNING *;

-- name: GetEntries :one

SELECT * FROM entries WHERE account_id = $1 LIMIT 1;
