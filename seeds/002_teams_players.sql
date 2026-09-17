-- Seed tim Arsenal & Chelsea beserta starting XI sesuai lineup di gambar referensi
-- (Arsenal 4-2-3-1 vs Chelsea 3-4-2-1). Aman dijalankan berulang: tim/pemain yang
-- sudah ada (dicek dari nama tim & nomor punggung per tim) tidak akan dibuat dobel.

-- === Tim ===

INSERT INTO teams (name, founded_year, address, city)
SELECT 'Arsenal', 1886, 'Emirates Stadium, Hornsey Rd', 'London'
WHERE NOT EXISTS (
    SELECT 1 FROM teams WHERE LOWER(name) = LOWER('Arsenal') AND deleted_at IS NULL
);

INSERT INTO teams (name, founded_year, address, city)
SELECT 'Chelsea', 1905, 'Stamford Bridge, Fulham Rd', 'London'
WHERE NOT EXISTS (
    SELECT 1 FROM teams WHERE LOWER(name) = LOWER('Chelsea') AND deleted_at IS NULL
);

-- === Pemain Arsenal ===
-- kolom: nama, tinggi(cm), berat(kg), posisi, nomor punggung

INSERT INTO players (team_id, name, height_cm, weight_kg, position, jersey_number)
SELECT t.id, p.name, p.height_cm, p.weight_kg, p.position, p.jersey_number
FROM teams t
JOIN (VALUES
    ('D. Raya',          183, 80, 'PENJAGA_GAWANG', 1),
    ('B. White',         180, 77, 'BERTAHAN',       4),
    ('E. Konsa',         180, 78, 'BERTAHAN',       15),
    ('G. Magalhães',     190, 87, 'BERTAHAN',       6),
    ('R. Calafiori',     186, 80, 'BERTAHAN',       33),
    ('D. Rice',          185, 82, 'GELANDANG',      41),
    ('M. Lewis-Skelly',  175, 68, 'GELANDANG',      49),
    ('B. Saka',          178, 70, 'PENYERANG',      7),
    ('M. Ødegaard',      178, 68, 'GELANDANG',      8),
    ('C. Tzolis',        178, 72, 'PENYERANG',      17),
    ('K. Havertz',       193, 82, 'PENYERANG',      29)
) AS p(name, height_cm, weight_kg, position, jersey_number) ON TRUE
WHERE t.name = 'Arsenal' AND t.deleted_at IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM players pl
      WHERE pl.team_id = t.id AND pl.jersey_number = p.jersey_number AND pl.deleted_at IS NULL
  );

-- === Pemain Chelsea ===

INSERT INTO players (team_id, name, height_cm, weight_kg, position, jersey_number)
SELECT t.id, p.name, p.height_cm, p.weight_kg, p.position, p.jersey_number
FROM teams t
JOIN (VALUES
    ('E. Martinez',      190, 85, 'PENJAGA_GAWANG', 1),
    ('J. Acheampong',    180, 75, 'BERTAHAN',       34),
    ('M. Lacroix',       190, 84, 'BERTAHAN',       5),
    ('W. Fofana',        186, 80, 'BERTAHAN',       3),
    ('J. Hato',          183, 76, 'BERTAHAN',       21),
    ('R. Lavia',         172, 65, 'GELANDANG',      45),
    ('R. James',         170, 68, 'GELANDANG',      24),
    ('P. Neto',          172, 65, 'BERTAHAN',       7),
    ('M. Rogers',        178, 72, 'PENYERANG',      17),
    ('C. Palmer',        189, 73, 'PENYERANG',      10),
    ('J. Pedro',         174, 68, 'PENYERANG',      9)
) AS p(name, height_cm, weight_kg, position, jersey_number) ON TRUE
WHERE t.name = 'Chelsea' AND t.deleted_at IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM players pl
      WHERE pl.team_id = t.id AND pl.jersey_number = p.jersey_number AND pl.deleted_at IS NULL
  );
