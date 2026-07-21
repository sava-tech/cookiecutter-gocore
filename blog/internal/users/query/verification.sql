
-- name: CreateVerification :one
INSERT INTO verifications (
    code, identifier_type, identifier, verification_type, expired_at, used, created_at
) VALUES (
    $1, $2, $3, $4, $5, false, NOW()
)
RETURNING *;

-- name: GetValidVerificationCode :one
SELECT * FROM verifications 
WHERE identifier = $1 
    AND code = $2 
    AND verification_type = $3 
    AND used = false 
    AND expired_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: GetVerificationHash :one
SELECT code FROM verifications 
WHERE identifier = $1 
    AND verification_type = $2 
    AND used = false
    AND expired_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkVerificationAsUsed :exec
UPDATE verifications 
SET used = true 
WHERE id = $1;

-- name: InvalidateVerificationCodes :exec
UPDATE verifications 
SET used = true 
WHERE identifier = $1 AND verification_type = $2 AND used = false;
