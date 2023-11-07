-- name: CreateAccount :one

INSERT INTO accounts ( owner, balance, currency ) VALUES ( $1, $2, $3 ) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1;

-- name: UpdateAccount :exec
UPDATE accounts set balance = $2 WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;