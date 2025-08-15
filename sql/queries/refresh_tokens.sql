-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, 
  created_at,
  updated_at,
  user_id,
  expires_at,
  revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT u.id FROM refresh_tokens r
LEFT JOIN users u ON u.id = r.user_id
WHERE r.token = $1
AND r.revoked_at IS NULL 
AND r.expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;
