CREATE TABLE match_events (
    id               BIGSERIAL PRIMARY KEY,
    uuid             UUID        NOT NULL DEFAULT gen_random_uuid(),
    match_id         BIGINT      NOT NULL REFERENCES matches(id),
    player_id        BIGINT      NOT NULL REFERENCES players(id),
    assist_player_id BIGINT      REFERENCES players(id),
    event_type       VARCHAR(20) NOT NULL
                     CHECK (event_type IN ('GOAL','YELLOW_CARD','RED_CARD','OWN_GOAL')),
    minute           SMALLINT    NOT NULL CHECK (minute BETWEEN 1 AND 130),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ,
    CONSTRAINT chk_assist_not_self CHECK (assist_player_id IS NULL OR assist_player_id <> player_id)
);
CREATE UNIQUE INDEX uq_match_events_uuid ON match_events (uuid);
CREATE INDEX idx_events_match ON match_events (match_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_events_match_type ON match_events (match_id, event_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_events_player ON match_events (player_id) WHERE deleted_at IS NULL;
