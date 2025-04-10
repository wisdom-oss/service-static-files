-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TYPE auth.scope_level AS ENUM('read', 'write', 'delete', '*');

CREATE TABLE IF NOT EXISTS
    auth.services (
        id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
        "name" TEXT NOT NULL UNIQUE,
        description TEXT,
        supported_scope_levels auth.scope_level[]
    );

INSERT INTO
    auth.services ("name", supported_scope_levels)
VALUES
    ('static-files', '{read, write, delete, *}');
-- +goose StatementEnd
