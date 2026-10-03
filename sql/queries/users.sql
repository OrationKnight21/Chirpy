-- name: CreateUser :one
INSERT INTO users(
    id, created_at, updated_at, email,hashed_password
) VALUES(
    $1,$2,$3,$4,$5
)
RETURNING *;
-- name: DeleteUsers :exec
DELETE FROM users;
-- name: GetUserByEmail :one
SELECT * FROM users where email = $1;
-- name: ChangeUserCreds :one
UPDATE users SET email = $1, hashed_password = $2, updated_at = NOW() where id = $3 RETURNING *;
-- name: ChangeChirpRedStatus :one
UPDATE users SET is_chirpy_red = TRUE, updated_at = NOW() where id = $1 RETURNING *;