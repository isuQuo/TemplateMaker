USE templatemaker;

CREATE TABLE users (
    id CHAR(36) NOT NULL PRIMARY KEY, 
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email); -- Add a unique constraint to the email column
