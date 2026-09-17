-- Seed user admin default untuk backend.
-- Password di-hash pakai pgcrypto (bcrypt cost 12), formatnya ($2a$...) kompatibel
-- dengan verifikasi golang.org/x/crypto/bcrypt di sisi Go — tidak perlu tool lain.
-- Aman dijalankan berulang: baris hanya ditambahkan kalau email belum ada.

INSERT INTO users (email, password)
SELECT 'admin@xyz.co.id', crypt('Admin#1234', gen_salt('bf', 12))
WHERE NOT EXISTS (
    SELECT 1 FROM users
    WHERE LOWER(email) = LOWER('admin@xyz.co.id') AND deleted_at IS NULL
);
