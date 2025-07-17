-- name: CreateConversationLog :one
INSERT INTO conversation_logs (
    user_id,
    model_id,
    request_payload,
    response_payload,
    prompt_tokens,
    completion_tokens,
    connection_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id, user_id, model_id, request_payload, response_payload, created_at, prompt_tokens, completion_tokens, connection_id;

-- name: ListConversationLogs :many
SELECT id, user_id, model_id, request_payload, response_payload, created_at, prompt_tokens, completion_tokens, connection_id
FROM conversation_logs
WHERE
    (sqlc.narg('user_id')::UUID IS NULL OR user_id = sqlc.narg('user_id')) AND
    (sqlc.narg('model_id')::UUID IS NULL OR model_id = sqlc.narg('model_id')) AND
    (sqlc.narg('connection_id')::UUID IS NULL OR connection_id = sqlc.narg('connection_id'))
ORDER BY created_at DESC
LIMIT sqlc.narg('limit')::BIGINT OFFSET sqlc.narg('offset')::BIGINT;

-- name: CountConversationLogs :one
SELECT COUNT(*)
FROM conversation_logs
WHERE
    (sqlc.narg('user_id')::UUID IS NULL OR user_id = sqlc.narg('user_id')) AND
    (sqlc.narg('model_id')::UUID IS NULL OR model_id = sqlc.narg('model_id')) AND
    (sqlc.narg('connection_id')::UUID IS NULL OR connection_id = sqlc.narg('connection_id'));

-- name: GetTotalTokensByProviderModelConnection :many
SELECT
    p.id AS provider_id,
    p.name AS provider_name,
    m.id AS model_id,
    m.proxy_model_id AS model_name,
    cl.connection_id,
    conn.name AS connection_name,
    SUM(cl.prompt_tokens + cl.completion_tokens) AS total_tokens
FROM
    conversation_logs cl
JOIN
    models m ON cl.model_id = m.id
JOIN
    connections conn ON m.connection_id = conn.id
JOIN
    providers p ON conn.provider_id::uuid = p.id
GROUP BY
    p.id,
    p.name,
    m.id,
    m.proxy_model_id,
    cl.connection_id,
    conn.name
ORDER BY
    p.id,
    m.id,
    cl.connection_id;

-- name: GetTotalPriceByProviderModelConnection :many
SELECT
    p.id AS provider_id,
    p.name AS provider_name,
    m.id AS model_id,
    m.proxy_model_id AS model_name,
    cl.connection_id,
    conn.name AS connection_name,
    SUM(
        (cl.prompt_tokens * m.price_input) +
        (cl.completion_tokens * m.price_output)
    )::NUMERIC AS total_price
FROM
    conversation_logs cl
JOIN
    models m ON cl.model_id = m.id
JOIN
    connections conn ON m.connection_id = conn.id
JOIN
    providers p ON conn.provider_id::uuid = p.id
GROUP BY
    p.id,
    p.name,
    m.id,
    m.proxy_model_id,
    cl.connection_id,
    conn.name
ORDER BY
    p.id,
    m.id,
    cl.connection_id;

-- name: GetTotalInputTokensByProviderModelConnection :many
SELECT
    p.id AS provider_id,
    p.name AS provider_name,
    m.id AS model_id,
    m.proxy_model_id AS model_name,
    cl.connection_id,
    conn.name AS connection_name,
    SUM(cl.prompt_tokens) AS total_input_tokens
FROM
    conversation_logs cl
JOIN
    models m ON cl.model_id = m.id
JOIN
    connections conn ON m.connection_id = conn.id
JOIN
    providers p ON conn.provider_id::uuid = p.id
GROUP BY
    p.id,
    p.name,
    m.id,
    m.proxy_model_id,
    cl.connection_id,
    conn.name
ORDER BY
    p.id,
    m.id,
    cl.connection_id;

-- name: GetTotalOutputTokensByProviderModelConnection :many
SELECT
    p.id AS provider_id,
    p.name AS provider_name,
    m.id AS model_id,
    m.proxy_model_id AS model_name,
    cl.connection_id,
    conn.name AS connection_name,
    SUM(cl.completion_tokens) AS total_output_tokens
FROM
    conversation_logs cl
JOIN
    models m ON cl.model_id = m.id
JOIN
    connections conn ON m.connection_id = conn.id
JOIN
    providers p ON conn.provider_id::uuid = p.id
GROUP BY
    p.id,
    p.name,
    m.id,
    m.proxy_model_id,
    cl.connection_id,
    conn.name
ORDER BY
    p.id,
    m.id,
    cl.connection_id;