# XYZ Football Team Management API

Backend REST API (JSON) untuk manajemen tim, pemain, jadwal, dan hasil pertandingan sepak bola amatir milik Perusahaan XYZ. Dibangun mengikuti [PRD](prd.md) dan [TRD](trd.md).

Stack: Go 1.26 + Gin + GORM + PostgreSQL 16 (TRD menyebut Go 1.22; dependency terbaru saat implementasi mensyaratkan minimum Go 1.26, sehingga versi toolchain disesuaikan).

## Menjalankan

```bash
cp .env.example .env
docker-compose up -d --build     # postgres + api
make migrate                     # jalankan migration
make seed                        # jalankan seed SQL (admin, tim, pemain, hasil pertandingan)
open http://localhost:8080/swagger/index.html
```

Login default (development): `admin@xyz.co.id` / `Admin#1234` — wajib diganti di environment non-development.

## Menjalankan tanpa Docker

```bash
cp .env.example .env
# sesuaikan DB_HOST/DB_PORT dst dengan instance PostgreSQL lokal
make migrate
make seed
make run
```

## Struktur Proyek

Mengikuti arsitektur layered (`handler → service → repository → entity`) sesuai TRD Bagian 3–4:

```
cmd/api        entry point
config/        load env & koneksi database
internal/
  entity/      model GORM
  dto/         request & response contract
  repository/  akses database
  service/     business rule
  handler/     HTTP layer
  middleware/  auth, rate limit, logging, recovery, cors
  router/      registrasi route
pkg/           response envelope, apperror, pagination, validator, jwt, uploader
migrations/    skema database (golang-migrate)
seeds/         data awal dalam SQL murni (admin, tim & pemain, hasil pertandingan)
```

## Data Awal (Seed)

Seed ditulis sebagai file `.sql` biasa di [seeds/](seeds/), dijalankan lewat `psql` — bukan kode Go — supaya isinya mudah dibaca dan diubah langsung:

- `001_admin.sql` — user admin default (`admin@xyz.co.id` / `Admin#1234`), password di-hash pakai `pgcrypto` (bcrypt) langsung di database.
- `002_teams_players.sql` — tim Arsenal & Chelsea beserta starting XI.
- `003_match_results.sql` — hasil pertandingan Arsenal 2-1 Chelsea beserta pencetak gol & menitnya.

Jalankan semua secara berurutan dengan:

```bash
make seed
```

Setiap file aman dijalankan berulang kali (idempoten) — baris yang sudah ada akan otomatis dilewati, tidak akan tercatat dobel.

## Testing

```bash
make test
```

## Dokumentasi API

Generate ulang Swagger setelah mengubah anotasi handler:

```bash
make swagger
```

Referensi lengkap kontrak API & aturan bisnis ada di [trd.md](trd.md) Bagian 6–8.
