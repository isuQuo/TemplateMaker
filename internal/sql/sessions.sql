USE templatemaker;

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY, -- 43 is the length of the SHA256 hash
    data BLOB NOT NULL, -- BLOB is a binary large object
    expiry TIMESTAMP(6) NOT NULL -- 6 is the number of decimal places for the fractional seconds
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);