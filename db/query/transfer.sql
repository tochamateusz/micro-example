-- name: CreateTransfer :one
INSERT INTO transfer (
"from_account_id", "to_account_id", "amount"
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransferById :one
SELECT * FROM transfer
WHERE id = $1 LIMIT 1;

-- name: GetTransfers :many
SELECT * FROM transfer
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetTransfersFromAccountId :many
SELECT * FROM transfer
WHERE from_account_id=$1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetTransfersToAccountId :many
SELECT * FROM transfer
WHERE to_account_id=$1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateTransfer :one
UPDATE transfer
SET amount=$2
WHERE id=$1
RETURNING *;
