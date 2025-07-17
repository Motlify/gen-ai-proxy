-- name: GetAPIKeyByHash :one
SELECT id, user_id, key_hash, name, created_at, last_used_at FROM api_keys
WHERE key_hash = $1;

-- name: GetAPIKey :one
SELECT id, user_id, key_hash, name, created_at, last_used_at FROM api_keys
WHERE key_hash = $1 AND user_id = $2;

-- name: CreateAPIKey :one
INSERT INTO api_keys (
    user_id,
    key_hash,
    name
) VALUES (
    $1, $2, $3
) RETURNING id, user_id, key_hash, name, created_at, last_used_at;

-- name: ListAPIKeys :many
SELECT id, user_id, name, created_at, last_used_at FROM api_keys
WHERE user_id = $1;

-- name: UpdateAPIKey :one
UPDATE api_keys
SET
    name = $2
WHERE name = $1 AND user_id = $3
RETURNING id, user_id, key_hash, name, created_at, last_used_at;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE api_keys
SET
    last_used_at = NOW()
WHERE id = $1;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys
WHERE id = $1 AND user_id = $2;