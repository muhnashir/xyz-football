CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    uuid       UUID         NOT NULL DEFAULT gen_random_uuid(),
    email      VARCHAR(150) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX uq_users_uuid  ON users (uuid);
CREATE UNIQUE INDEX uq_users_email_active
    ON users (LOWER(email)) WHERE deleted_at IS NULL;
