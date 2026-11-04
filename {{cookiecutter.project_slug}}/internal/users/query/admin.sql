-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*)
FROM users
WHERE is_active = true;

-- name: ListRegisteredUsers :many
SELECT id, email, first_name, last_name, email_verified, created_at
FROM users
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
