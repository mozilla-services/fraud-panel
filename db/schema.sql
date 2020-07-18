CREATE DATABASE fpnl;
CREATE ROLE fpnlapi;
ALTER ROLE fpnlapi WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB LOGIN PASSWORD 'fpnlapi';
\c fpnl
CREATE TABLE users (
      id            SERIAL PRIMARY KEY,
      name          VARCHAR NOT NULL,
      created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
      last_modified TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
      verified      BOOLEAN NOT NULL
);
GRANT SELECT, INSERT, UPDATE ON users TO fpnlapi;
GRANT USAGE ON users_id_seq TO fpnlapi;