-- name: CreateConnection :one
INSERT INTO connections (
    user_id,
    provider_id,
    encrypted_api_key,
    name
) VALUES (
    $1, $2, $3, $4
) RETURNING id, user_id, provider_id, encrypted_api_key, name, created_at;

-- name: GetConnection :one
SELECT c.id, c.user_id, c.provider_id, c.encrypted_api_key, c.name, c.created_at, p.type as provider_type, c.deleted_at
FROM connections c
JOIN providers p ON p.id = c.provider_id::uuid
WHERE c.id = $1 AND c.user_id = $2 AND c.deleted_at IS NULL;

-- name: GetConnectionByProvider :one
SELECT id, user_id, provider_id, encrypted_api_key, name, created_at, deleted_at FROM connections
WHERE user_id = $1 AND provider_id = $2 AND deleted_at IS NULL;

-- name: ListConnections :many
SELECT id, user_id, provider_id, name, created_at, deleted_at FROM connections
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: ListConnectionsByProviderID :many
SELECT id, user_id, provider_id, encrypted_api_key, name, created_at, deleted_at FROM connections
WHERE provider_id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: SoftDeleteConnection :exec
UPDATE connections
SET deleted_at = NOW()
WHERE id = $1 AND user_id = $2;
