-- name: GetModelByProxyModelID :one
SELECT * FROM models WHERE proxy_model_id = $1 AND user_id = $2;

-- name: CreateModel :one
INSERT INTO models (
    id,
    user_id,
    connection_id,
    proxy_model_id,
    provider_model_id,
    thinking,
    tools_usage,
    price_input,
    price_output
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetModel :one
SELECT * FROM models WHERE id = $1 AND user_id = $2;

-- name: ListModels :many
SELECT * FROM models WHERE user_id = $1 AND deleted_at IS NULL;

-- name: UpdateModel :one
UPDATE models
SET
    proxy_model_id = $3,
    provider_model_id = $4,
    thinking = $5,
    tools_usage = $6,
    price_input = $7,
    price_output = $8
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: SoftDeleteModel :exec
UPDATE models
SET deleted_at = NOW()
WHERE id = $1 AND user_id = $2;
