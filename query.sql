-- name: GetUserById :one
SELECT * FROM users WHERE id = $1

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1

-- name: GetUserByOAuthClientID :one
SELECT * FROM users WHERE id = (SELECT user_id FROM oauth_clients WHERE client_id = $1)

-- name: GetUserGroupsByUserID :many
SELECT * FROM user_groups WHERE id IN (SELECT user_group_id FROM user_group_map WHERE user_id = $1)
