CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE user_groups (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE user_group_map (
    user_id INTEGER NOT NULL REFERENCES users(id),
    user_group_id INTEGER NOT NULL REFERENCES user_groups(id),
    PRIMARY KEY (user_id, user_group_id)
);

CREATE TABLE oauth_clients (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR(100) NOT NULL UNIQUE,
    client_secret VARCHAR(100) NOT NULL,
    redirect_uri VARCHAR(255) NOT NULL,
    grant_types VARCHAR(255) NOT NULL,
    scope VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE oauth_access_tokens (
    access_token VARCHAR(100) NOT NULL PRIMARY KEY,
    client_id VARCHAR(100) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    expires TIMESTAMP NOT NULL,
    scope VARCHAR(255) NOT NULL
);

CREATE TABLE oauth_refresh_tokens (
    refresh_token VARCHAR(100) NOT NULL PRIMARY KEY,
    client_id VARCHAR(100) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    expires TIMESTAMP NOT NULL,
    scope VARCHAR(255) NOT NULL
);

CREATE TABLE oauth_auth_codes (
    auth_code VARCHAR(100) NOT NULL PRIMARY KEY,
    client_id VARCHAR(100) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    redirect_uri VARCHAR(255) NOT NULL,
    expires TIMESTAMP NOT NULL,
    scope VARCHAR(255) NOT NULL
);

CREATE TABLE oauth_scopes (
    scope TEXT NOT NULL,
    is_default BOOLEAN
);

CREATE TABLE oauth_jwt (
    client_id VARCHAR(100) NOT NULL PRIMARY KEY,
    subject VARCHAR(100),
    public_key VARCHAR(2000)
);

CREATE TABLE oauth_public_keys (
    client_id VARCHAR(100) NOT NULL PRIMARY KEY,
    public_key VARCHAR(2000) NOT NULL,
    private_key VARCHAR(2000) NOT NULL,
    encryption_algorithm VARCHAR(100) NOT NULL DEFAULT 'RS256'
);
