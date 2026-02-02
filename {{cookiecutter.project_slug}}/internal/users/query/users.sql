-- name: CreateUser :one
INSERT INTO users (
  email,
  username,
  phone_number,
  avatar,
  age,
  gender,
  hashed_password
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users     
WHERE email = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;


-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: UpdateUser :one
UPDATE users
SET
  email = COALESCE($2, email),
  username = COALESCE($3, username),
  phone_number = COALESCE($4, phone_number),
  avatar = COALESCE($5, avatar),
  age = COALESCE($6, age),
  gender = COALESCE($7, gender),
  is_active = COALESCE($8, is_active),
  updated_at = now()
WHERE id = $1
RETURNING *;


-- name: UpdatePassword :exec
UPDATE users
SET
  hashed_password = $2,
  password_changed_at = now(),
  updated_at = now()
WHERE id = $1;


-- name: DeactivateUser :exec
UPDATE users
SET
  is_active = false,
  updated_at = now()
WHERE id = $1;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
