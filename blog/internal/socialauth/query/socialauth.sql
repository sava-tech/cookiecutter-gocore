-- name: CreateSocialUser :one
INSERT INTO social_auth_users (
    id, user_id, provider, provider_id, email, name, avatar_url, access_token
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetSocialUserByProviderID :one
SELECT * FROM social_auth_users 
WHERE provider = $1 AND provider_id = $2 
LIMIT 1;

-- name: GetSocialUserByUserID :one
SELECT * FROM social_auth_users 
WHERE user_id = $1 
LIMIT 1;

-- name: UpdateSocialUser :one
UPDATE social_auth_users 
SET 
    access_token = COALESCE($2, access_token),
    avatar_url = COALESCE($3, avatar_url),
    updated_at = NOW()
WHERE id = $1
RETURNING *;