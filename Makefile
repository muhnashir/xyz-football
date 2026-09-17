include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

run:      ; go run ./cmd/api
migrate:  ; migrate -path migrations -database "$(DB_URL)" up
rollback: ; migrate -path migrations -database "$(DB_URL)" down 1
seed:     ; go run ./cmd/seeder
swagger:  ; swag init -g cmd/api/main.go -o docs
test:     ; go test ./... -cover
up:       ; docker-compose up -d --build
down:     ; docker-compose down

.PHONY: run migrate rollback seed swagger test up down
