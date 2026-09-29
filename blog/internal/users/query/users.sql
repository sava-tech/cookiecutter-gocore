-- name: CreateUser :one
INSERT INTO users (
  email,
  first_name,
  last_name,
  full_name,
  phone_number,
  avatar,
  age,
  gender,
  email_verified,
  account_type,
  hashed_password,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
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

-- name: GetUserByFirstName :one
SELECT *
FROM users
WHERE first_name = $1
LIMIT 1;


-- name: UpdateUser :one
UPDATE users
SET
  email = COALESCE($2, email),
  first_name = COALESCE($3, first_name),
  last_name = COALESCE($4, last_name),
  full_name = COALESCE($5, full_name),
  phone_number = COALESCE($6, phone_number),
  avatar = COALESCE($7, avatar),
  age = COALESCE($8, age),
  gender = COALESCE($9, gender),
  is_active = COALESCE($10, is_active),
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: VerifyUserEmail :exec
UPDATE users 
SET email_verified = true, updated_at = NOW() 
WHERE email = $1;

-- name: UpdatePassword :exec
UPDATE users
SET
  hashed_password = $2,
  password_changed_at = now(),
  updated_at = now()
WHERE id = $1;

-- name: GetUserByPhoneNumber :one
SELECT * FROM users WHERE phone_number = $1 LIMIT 1;

-- name: DeactivateUser :exec
UPDATE users
SET
  is_active = false,
  updated_at = now()
WHERE id = $1;


-- name: UpdateUserPassword :exec
UPDATE users 
SET hashed_password = $2, updated_at = NOW() 
WHERE email = $1;


