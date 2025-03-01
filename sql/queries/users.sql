-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserFromID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUserEmailPW :one
UPDATE users
SET updated_at = NOW(), email = $2, hashed_password = $3
WHERE id = $1
RETURNING *;

-- name: UpgradeUserRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;