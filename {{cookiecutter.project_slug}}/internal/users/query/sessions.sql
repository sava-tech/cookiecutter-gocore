-- name: CreateSession :one
INSERT INTO sessions (
    id, user_id, email, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, NOW()
)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions WHERE id = $1 LIMIT 1;

-- name: UpdateSession :one
UPDATE sessions 
SET is_blocked = $2, refresh_token = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;