CREATE TABLE players (
    id            BIGSERIAL PRIMARY KEY,
    uuid          UUID         NOT NULL DEFAULT gen_random_uuid(),
    team_id       BIGINT       NOT NULL REFERENCES teams(id),
    name          VARCHAR(100) NOT NULL,
    height_cm     SMALLINT     NOT NULL CHECK (height_cm BETWEEN 100 AND 250),
    weight_kg     SMALLINT     NOT NULL CHECK (weight_kg BETWEEN 30 AND 200),
    position      VARCHAR(20)  NOT NULL
                  CHECK (position IN ('PENYERANG','GELANDANG','BERTAHAN','PENJAGA_GAWANG')),
    jersey_number SMALLINT     NOT NULL CHECK (jersey_number BETWEEN 1 AND 99),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);
CREATE UNIQUE INDEX uq_players_uuid ON players (uuid);
CREATE UNIQUE INDEX uq_players_team_jersey_active
    ON players (team_id, jersey_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_players_team ON players (team_id) WHERE deleted_at IS NULL;
