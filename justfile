# Set up variables
BINARY := "lbox"
SRC := "./main.go"
BUILD_DIR := "./bin"
MIGRATION_DIR := "./db/migrations"
DB_DIR := "./app.db"

run:
	sqlc generate
	templ generate
	go run {{SRC}}

generate:
	sqlc generate
	templ generate

letterboxd:
	go run ./cmd/letterboxd/main.go

film:
	go run ./cmd/film/main.go


build:
    mkdir -p {{BUILD_DIR}}
    go build -o {{BUILD_DIR}}/{{BINARY}} {{SRC}}

clean:
    rm -rf {{BUILD_DIR}}

reset:
	-rm {{DB_DIR}}
	just db-create
	just migrate-up


# Create Database
[group('db')]
db-create:
	touch {{DB_DIR}}

# Deletes the DB giving you a choice.
[group('db')]
db-delete:
	@read -p "Do you want to delete the DB (you'll loose all data)? [y/n] " choice; \
	if [ "$$choice" != "y" ] && [ "$$choice" != "Y" ]; then \
		echo "Exiting..."; \
		exit 1; \
	else \
		rm -f db/app.db; \
	fi; \

[group('db')]
[group('migration')]
migration-create arg_name:
	@ mkdir -p db/migrations
	goose -dir {{MIGRATION_DIR}} create {{arg_name}} sql

[group('db')]
[group('migration')]
migrate-up:
    GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING={{DB_DIR}} goose up -dir {{MIGRATION_DIR}}
