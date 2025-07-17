-- name: CreateUser :one
INSERT INTO users (
    username,
    password_hash
) VALUES (
    $1, $2
) RETURNING id, username, password_hash, created_at;

-- name: GetUserPasswordHash :one
SELECT
    password_hash as hash_value
FROM
    users
WHERE
    username = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;