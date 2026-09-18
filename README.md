# XYZ Football Team Management API

Backend REST API (JSON) untuk manajemen tim, pemain, jadwal, dan hasil pertandingan sepak bola amatir milik Perusahaan XYZ.

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
make swagger                     # wajib sekali di awal clone — lihat catatan di bawah
make migrate
make seed
make run
```

> `docs/` (hasil generate Swagger) sengaja tidak di-commit ke git (lihat `.gitignore`), padahal `cmd/api/main.go` meng-import package tersebut — jadi tanpa `make swagger` lebih dulu, `make run`/`go build` akan gagal dengan error `no required module provides package .../docs`. Jalur Docker (`docker-compose up`) sudah otomatis menjalankan generate ini di dalam build stage, jadi tidak perlu langkah manual tambahan.

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

### Swagger UI (cara utama menjalankan/mencoba API)

Setelah server jalan (`make run` atau `docker-compose up`), buka:

```
http://localhost:8080/swagger/index.html
```

Semua endpoint (Auth, Teams, Players, Matches, Reports) terdaftar lengkap dengan contoh payload, response, dan kode error. Klik **Authorize**, isi `Bearer <access_token>` hasil `POST /auth/login`, lalu setiap endpoint yang butuh login bisa langsung dicoba lewat tombol **Try it out** — tidak perlu tool tambahan.

Generate ulang Swagger setelah mengubah anotasi handler:

```bash
make swagger
```

> Perintah di atas sengaja memakai `go run github.com/swaggo/swag/cmd/swag@v1.8.12` (bukan `swag` global), supaya versi generator selalu cocok dengan `github.com/swaggo/swag` di `go.mod`. Menjalankan `swag init` memakai CLI global bisa gagal build (`unknown field LeftDelim/RightDelim ...`) kalau versinya lebih baru dari yang dipakai project.

### Import ke Postman (opsional)

Swagger spec di atas adalah OpenAPI 2.0 murni, jadi bisa langsung diimpor ke Postman tanpa perlu file collection terpisah yang harus dirawat manual:

1. Buka Postman → **Import**.
2. Pilih **Link**, masukkan `http://localhost:8080/swagger/doc.json` (server harus sedang jalan), atau **File** lalu pilih [docs/swagger.json](docs/swagger.json) dari repo ini.
3. Postman otomatis membuat collection berisi seluruh endpoint beserta skema request/response.
4. Set environment variable `baseUrl = http://localhost:8080/api/v1`, lalu simpan `access_token` dari response `POST /auth/login` ke variable (mis. `token`) dan pakai di header `Authorization: Bearer {{token}}` untuk endpoint yang protected.

## Jawaban Soal Rekrutmen

Selain source code ini (jawaban soal nomor 1 — implementasi API), dua soal lain dijawab dalam file teks terpisah di root repo:

- [soal_dan_jawaban_nomor_2](soal_dan_jawaban_nomor_2) — evaluasi produk AYO (web `ayo.co.id` & aplikasi mobile iOS/Android): temuan bug, celah validasi/keamanan, dan rekomendasi peningkatan UX.
- [soal_dan_jawaban_nomor_3](soal_dan_jawaban_nomor_3) — jawaban pertanyaan rekrutmen (motivasi, keunggulan, pengalaman relevan, dan rencana 3 bulan pertama).
