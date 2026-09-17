-- Seed hasil pertandingan sesuai screenshot skor akhir referensi:
--   Arsenal 2 - 1 Chelsea, Sun 6 Sept 2026
--   Gol: M. Rogers 2' (Chelsea), K. Havertz 25' (Arsenal), M. Ødegaard 50' (Arsenal)
-- Harus dijalankan setelah 002_teams_players.sql.
--
-- Idempoten: CTE new_match hanya INSERT kalau pertandingan (home, away, tanggal)
-- belum ada. Kalau sudah ada, new_match tidak menghasilkan baris, sehingga INSERT
-- match_events di bawahnya otomatis ikut ter-skip (tidak ada gol dobel).

WITH new_match AS (
    INSERT INTO matches (match_date, match_time, home_team_id, away_team_id, status, home_score, away_score)
    SELECT '2026-09-06', '14:00', home.id, away.id, 'finished', 2, 1
    FROM teams home, teams away
    WHERE home.name = 'Arsenal' AND home.deleted_at IS NULL
      AND away.name = 'Chelsea' AND away.deleted_at IS NULL
      AND NOT EXISTS (
          SELECT 1 FROM matches m
          WHERE m.home_team_id = home.id
            AND m.away_team_id = away.id
            AND m.match_date = '2026-09-06'
            AND m.deleted_at IS NULL
      )
    RETURNING id
)
INSERT INTO match_events (match_id, player_id, event_type, minute)
SELECT nm.id, pl.id, 'GOAL', g.minute
FROM new_match nm
JOIN (VALUES
    ('Chelsea', 17, 2),   -- M. Rogers
    ('Arsenal', 29, 25),  -- K. Havertz
    ('Arsenal', 8,  50)   -- M. Ødegaard
) AS g(team_name, jersey_number, minute) ON TRUE
JOIN teams t   ON t.name = g.team_name AND t.deleted_at IS NULL
JOIN players pl ON pl.team_id = t.id AND pl.jersey_number = g.jersey_number AND pl.deleted_at IS NULL;
