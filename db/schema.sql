CREATE DATABASE fpnl;
CREATE ROLE fpnlapi;
ALTER ROLE fpnlapi WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB LOGIN PASSWORD 'fpnlapi';
\c fpnl
CREATE EXTENSION pgcrypto;
CREATE TABLE account (
      id            SERIAL PRIMARY KEY,
      created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
      last_modified TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
      active        BOOLEAN NOT NULL
);
CREATE TABLE member (
      id            SERIAL PRIMARY KEY,
      name          VARCHAR NOT NULL,
      email         VARCHAR UNIQUE NOT NULL,
      created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
      last_modified TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE TABLE session (
      id            CHAR(64) PRIMARY KEY DEFAULT encode(sha256(gen_random_bytes(128)), 'hex'),
      member_id     SERIAL REFERENCES member(id) NOT NULL,
      created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE TABLE service (
      id            SERIAL PRIMARY KEY,
      account_id    SERIAL REFERENCES account(id) NOT NULL,
      created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
      last_modified TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE TABLE config (
      id            SERIAL PRIMARY KEY,
      service_id    SERIAL REFERENCES service(id) NOT NULL,
      created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
      last_modified TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
      revision      VARCHAR NOT NULL
);
CREATE TABLE account_permission (
      account_id  SERIAL REFERENCES account(id),
      member_id   SERIAL REFERENCES member(id),
      role        VARCHAR NOT NULL
);
CREATE TABLE service_permission (
      service_id  SERIAL REFERENCES service(id),
      member_id   SERIAL REFERENCES member(id),
      role        VARCHAR NOT NULL
);
CREATE TABLE config_permission (
      config_id  SERIAL REFERENCES config(id),
      member_id  SERIAL REFERENCES member(id),
      role       VARCHAR NOT NULL
);
GRANT SELECT, INSERT, UPDATE ON account, member, session, service, config, account_permission, service_permission, config_permission TO fpnlapi;
GRANT USAGE ON account_id_seq, member_id_seq, service_id_seq, config_id_seq TO fpnlapi;