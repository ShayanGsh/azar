-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: AddUser :exec
INSERT INTO users (username, email, password) VALUES ($1, $2, $3);

-- name: UpdatePassword :exec
UPDATE users SET password = $1 WHERE id = $2;

-- name: UpdateUser :exec
UPDATE users SET username = $1, email = $2, updated_at = $3 WHERE id = $4;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: DeleteUserGroup :exec
DELETE FROM user_groups WHERE id = $1;

-- name: VerifyUser :exec
SELECT * FROM users WHERE username = $1 OR email = $1 AND password = $2;

-- name: VerifyUserByEmail :exec
SELECT * FROM users WHERE email = $1 AND password = $2;