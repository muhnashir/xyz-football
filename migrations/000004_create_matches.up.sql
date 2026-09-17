CREATE TABLE matches (
    id           BIGSERIAL PRIMARY KEY,
    uuid         UUID        NOT NULL DEFAULT gen_random_uuid(),
    match_date   DATE        NOT NULL,
    match_time   TIME        NOT NULL,
    home_team_id BIGINT      NOT NULL REFERENCES teams(id),
    away_team_id BIGINT      NOT NULL REFERENCES teams(id),
    status       VARCHAR(20) NOT NULL DEFAULT 'scheduled'
                 CHECK (status IN ('scheduled','in_progress','finished','postpone')),
    home_score   SMALLINT    NOT NULL DEFAULT 0 CHECK (home_score >= 0),
    away_score   SMALLINT    NOT NULL DEFAULT 0 CHECK (away_score >= 0),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT chk_home_not_away CHECK (home_team_id <> away_team_id)
);
CREATE UNIQUE INDEX uq_matches_uuid ON matches (uuid);
CREATE INDEX idx_matches_schedule ON matches (match_date, match_time, id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_home ON matches (home_team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_away ON matches (away_team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_status ON matches (status) WHERE deleted_at IS NULL;
