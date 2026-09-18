include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

run:      ; go run ./cmd/api
migrate:  ; migrate -path migrations -database "$(DB_URL)" up
rollback: ; migrate -path migrations -database "$(DB_URL)" down 1
swagger:  ; go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/api/main.go -o docs
test:     ; go test ./... -cover
up:       ; docker-compose up -d --build
down:     ; docker-compose down

# Jalankan file seed SQL di folder seeds/ berurutan sesuai nomor prefix-nya.
# Aman dijalankan berulang kali (idempoten).
seed:
	@for f in seeds/*.sql; do \
		echo "==> $$f"; \
		psql "$(DB_URL)" -v ON_ERROR_STOP=1 -f $$f || exit 1; \
	done

.PHONY: run migrate rollback swagger test up down seed
