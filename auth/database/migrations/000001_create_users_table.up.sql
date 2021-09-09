CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id text NOT NULL PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    first_name text NOT NULL,
    last_name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password bytea NOT NULL,
    activated bool NOT NULL
);

