-- +migrate Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Up
CREATE TABLE user_groups (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Up
CREATE TABLE user_group_map (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_group_id INTEGER NOT NULL REFERENCES user_groups(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, user_group_id)
);

-- +migrate Up
CREATE TABLE oauth_clients (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR(100) NOT NULL UNIQUE,
    client_secret VARCHAR(100) NOT NULL,
    redirect_uri VARCHAR(255) NOT NULL,
    grant_types VARCHAR(255) NOT NULL,
    scope VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Up
CREATE TABLE oauth_access_tokens (
    access_token VARCHAR(100) NOT NULL PRIMARY KEY,
    client_id VARCHAR(100) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    expires TIMESTAMP NOT NULL,
    scope VARCHAR(255) NOT NULL
);

-- +migrate Up
CREATE TABLE oauth_refresh_tokens (
    refresh_token VARCHAR(100) NOT NULL PRIMARY KEY,
    client_id VARCHAR(100) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    expires TIMESTAMP NOT NULL,
    scope VARCHAR(255) NOT NULL
);

-- +migrate Up
CREATE TABLE oauth_auth_codes (
    auth_code VARCHAR(100) NOT NULL PRIMARY KEY,
    client_id VARCHAR(100) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    redirect_uri VARCHAR(255) NOT NULL,
    expires TIMESTAMP NOT NULL,
    scope VARCHAR(255) NOT NULL
);

-- +migrate Up
CREATE TABLE oauth_scopes (
    scope TEXT NOT NULL,
    is_default BOOLEAN
);

-- +migrate Up
CREATE TABLE oauth_jwt (
    client_id VARCHAR(100) NOT NULL PRIMARY KEY,
    subject VARCHAR(100),
    public_key VARCHAR(2000)
);

-- +migrate Up
CREATE TABLE oauth_public_keys (
    client_id VARCHAR(100) NOT NULL PRIMARY KEY,
    public_key VARCHAR(2000) NOT NULL,
    private_key VARCHAR(2000) NOT NULL,
    encryption_algorithm VARCHAR(100) NOT NULL DEFAULT 'RS256'
);

-- +migrate Up
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Up
CREATE TABLE permission_groups (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Up
CREATE TABLE permission_group_map (
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    permission_group_id INTEGER NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE,
    PRIMARY KEY (permission_id, permission_group_id)
);

-- +migrate Up
CREATE TABLE permission_group_user_map (
    permission_group_id INTEGER NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (permission_group_id, user_id)
);

-- +migrate Up
CREATE TABLE permission_group_user_group_map (
    permission_group_id INTEGER NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE,
    user_group_id INTEGER NOT NULL REFERENCES user_groups(id) ON DELETE CASCADE,
    PRIMARY KEY (permission_group_id, user_group_id)
);


-- +migrate Down
DROP TABLE users;

-- +migrate Down
DROP TABLE user_groups;

-- +migrate Down
DROP TABLE user_group_map;

-- +migrate Down
DROP TABLE oauth_clients;

-- +migrate Down
DROP TABLE oauth_access_tokens;

-- +migrate Down
DROP TABLE oauth_refresh_tokens;

-- +migrate Down
DROP TABLE oauth_auth_codes;

-- +migrate Down
DROP TABLE oauth_scopes;

-- +migrate Down
DROP TABLE oauth_jwt;

-- +migrate Down
DROP TABLE oauth_public_keys;

-- +migrate Down
DROP TABLE permissions;

-- +migrate Down
DROP TABLE permission_groups;

-- +migrate Down
DROP TABLE permission_group_map;

-- +migrate Down
DROP TABLE permission_group_user_map;

-- +migrate Down
DROP TABLE permission_group_user_group_map;
