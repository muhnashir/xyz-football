CREATE TABLE teams (
    id           BIGSERIAL PRIMARY KEY,
    uuid         UUID         NOT NULL DEFAULT gen_random_uuid(),
    name         VARCHAR(100) NOT NULL,
    founded_year SMALLINT     NOT NULL CHECK (founded_year BETWEEN 1800 AND 2100),
    logo_url     VARCHAR(255),
    address      TEXT         NOT NULL,
    city         VARCHAR(100) NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);
CREATE UNIQUE INDEX uq_teams_uuid ON teams (uuid);
CREATE UNIQUE INDEX uq_teams_name_active
    ON teams (LOWER(name)) WHERE deleted_at IS NULL;
CREATE INDEX idx_teams_city ON teams (city) WHERE deleted_at IS NULL;
