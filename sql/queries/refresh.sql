-- name: CreateToken :one
INSERT INTO refresh_tokens(
    token,created_at,updated_at,user_id,expires_at,revoked_at
) VALUES(
    $1,$2,$3,$4,$5,$6
) 
RETURNING *;
-- name: GetUserRefreshToken :one
SELECT user_id from refresh_tokens where token = $1 AND revoked_at IS NULL AND expires_at>NOW();
-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW(), updated_at = NOW() where token = $1;