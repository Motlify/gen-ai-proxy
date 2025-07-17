-- name: CreateProvider :one
INSERT INTO providers (
    id,
    user_id,
    name,
    base_url,
    type
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, user_id, name, base_url, type;

-- name: GetProvider :one
SELECT id, user_id, name, base_url, type, deleted_at FROM providers WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: ListProviders :many
SELECT id, user_id, name, base_url, type, deleted_at FROM providers WHERE user_id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteProvider :exec
UPDATE providers
SET deleted_at = NOW()
WHERE id = $1 AND user_id = $2;