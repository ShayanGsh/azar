-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByOAuthClientID :one
SELECT * FROM users WHERE id = (SELECT user_id FROM oauth_clients WHERE client_id = $1);

-- name: GetUserGroupsByUserID :many
SELECT * FROM user_groups WHERE id IN (SELECT user_group_id FROM user_group_map WHERE user_id = $1);

-- name: AddUser :exec
INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);

-- name: AddUserGroup :exec
INSERT INTO user_groups (group_name, created_at, updated_at) VALUES ($1, $2, $3);

-- name: UpdatePassword :exec
UPDATE users SET password = $1, updated_at = $2 WHERE id = $3;

-- name: UpdateUser :exec
UPDATE users SET username = $1, email = $2, updated_at = $3 WHERE id = $4;

-- name: AddUserGroupMap :exec
INSERT INTO user_group_map (user_id, user_group_id) VALUES ($1, $2);

-- name: RemoveUserGroupMap :exec
DELETE FROM user_group_map WHERE user_id = $1 AND user_group_id = $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: DeleteUserGroup :exec
DELETE FROM user_groups WHERE id = $1;

-- name: DeleteUserGroupMap :exec
DELETE FROM user_group_map WHERE user_group_id = $1;

-- name: DeleteUserGroupMapByUserID :exec
DELETE FROM user_group_map WHERE user_id = $1;

-- name: AddPermission :exec
INSERT INTO permissions (name, created_at, updated_at) VALUES ($1, $2, $3);

-- name: AddPermissionGroup :exec
INSERT INTO permission_groups (group_name, created_at, updated_at) VALUES ($1, $2, $3);

-- name: AddPermissionGroupMap :exec
INSERT INTO permission_group_map (permission_id, permission_group_id) VALUES ($1, $2);

-- name: UpdatePermission :exec
UPDATE permissions SET name = $1, updated_at = $2 WHERE id = $3;

-- name: UpdatePermissionGroup :exec
UPDATE permission_groups SET group_name = $1, updated_at = $2 WHERE id = $3;

-- name: RemovePermissionGroupMap :exec
DELETE FROM permission_group_map WHERE permission_id = $1 AND permission_group_id = $2;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE id = $1;

-- name: DeletePermissionGroup :exec
DELETE FROM permission_groups WHERE id = $1;

-- name: DeletePermissionGroupMap :exec
DELETE FROM permission_group_map WHERE permission_group_id = $1;

-- name: VerifyUser :exec
SELECT * FROM users WHERE username = $1 OR email = $1 AND password = $2;

-- name: VerifyUserByEmail :exec
SELECT * FROM users WHERE email = $1 AND password = $2;