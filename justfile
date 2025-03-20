MIGRATION_DIR := "./migrations"
DB_STRING := "postgres://postgres@localhost:5432/app?sslmode=disable"

run:
    sqlc generate
    templ generate
    go run main.go

generate:
    sqlc generate
    templ generate

letterboxd:
    go run ./cmd/letterboxd/main.go

film:
    go run ./cmd/film/main.go


reset:
	docker compose down
	docker compose up -d
	sleep 2
	just migrate-up

[group('migration')]
migration-create arg_name:
	@mkdir -p {{MIGRATION_DIR}}
	goose -dir {{MIGRATION_DIR}} create {{arg_name}} sql

[group('migration')]
migrate-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="{{DB_STRING}}" goose up -dir {{MIGRATION_DIR}}

[group('migration')]
migrate-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="{{DB_STRING}}" goose down -dir {{MIGRATION_DIR}}

