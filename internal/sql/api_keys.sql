CREATE TABLE api_keys (
    name VARCHAR(255) NOT NULL UNIQUE,
    api_key VARCHAR(255) NOT NULL UNIQUE
);