-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS auth;

DO $$ BEGIN
CREATE TYPE auth.scope_level AS ENUM('read', 'write', 'delete', '*');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

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
