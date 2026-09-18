FROM golang:1.26-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# docs/ (Swagger) tidak di-commit ke git (lihat .gitignore) dan digenerate di sini
# supaya image selalu bisa dibangun dari clone bersih, tanpa bergantung pada
# state lokal developer. Versi swag dikunci agar sama dengan go.mod.
RUN go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/api/main.go -o docs
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/api ./api
COPY --from=builder /app/docs ./docs
COPY migrations ./migrations
COPY seeds ./seeds
EXPOSE 8080
CMD ["./api"]
