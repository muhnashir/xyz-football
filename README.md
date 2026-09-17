# XYZ Football Team Management API

Backend REST API (JSON) untuk manajemen tim, pemain, jadwal, dan hasil pertandingan sepak bola amatir milik Perusahaan XYZ. Dibangun mengikuti [PRD](prd.md) dan [TRD](trd.md).

Stack: Go 1.26 + Gin + GORM + PostgreSQL 16 (TRD menyebut Go 1.22; dependency terbaru saat implementasi mensyaratkan minimum Go 1.26, sehingga versi toolchain disesuaikan).

## Menjalankan

```bash
cp .env.example .env
docker-compose up -d --build     # postgres + api
make migrate                     # jalankan migration
make seed                        # buat user admin default
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
cmd/seeder     seed user admin
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
```

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
