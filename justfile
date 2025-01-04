# Set up variables
binary := "lbox"
src := "./main.go"
build-dir := "./bin"
migration-dir := "./db/migrations"
db-dir := "./app.db"
db-type := "postgres" # "sqlite" or "postgres"

# PostgreSQL configurations
pg-user := "arar"                       # PostgreSQL username
pg-password := "password"               # PostgreSQL password
pg-dbname := "arar"                    # PostgreSQL database name
pg-host := "localhost"                 # PostgreSQL host
pg-sslmode := "disable"                # PostgreSQL SSL mode
db-string := "postgres://" + pg-user + "@" + pg-host + ":5432/" + pg-dbname + "?sslmode=" + pg-sslmode

# Default recipe (optional)
default:
    just --list

run:
    sqlc generate
    templ generate
    go run {{src}}

generate:
    sqlc generate
    templ generate

letterboxd:
    go run ./cmd/letterboxd/main.go

film:
    go run ./cmd/film/main.go

build:
    mkdir -p {{build-dir}}
    go build -o {{build-dir}}/{{binary}} {{src}}

clean:
    rm -rf {{build-dir}}

reset: db-delete db-create migrate-up

# Create Database
[group('db')]
db-create:
    #!/usr/bin/env sh
    if [ "{{db-type}}" = "sqlite" ]; then \
        touch {{db-dir}}; \
    else \
        echo "Creating database handled by direnv..."; \
        # psql -c "CREATE DATABASE {{pg-dbname}};" -U {{pg-user}}; \
    fi

# Deletes all tables in the DB giving you a choice.
[group('db')]
db-delete:
    #!/usr/bin/env sh
        if [ "{{db-type}}" = "sqlite" ]; then \
            rm -f {{db-dir}}; \
        else \
            echo "Dropping all tables in the PostgreSQL database..."; \
            psql -U {{pg-user}} -d {{pg-dbname}} -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"; \
        fi; \

[group('db')]
[group('migration')]
migration-create name:
    mkdir -p db/migrations
    goose -dir {{migration-dir}} create {{name}} sql

[group('db')]
[group('migration')]
migrate-up:
    #!/usr/bin/env sh
    if [ "{{db-type}}" = "sqlite" ]; then \
        GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING={{db-dir}} goose up -dir {{migration-dir}}; \
    else \
        GOOSE_DRIVER=postgres GOOSE_DBSTRING="{{db-string}}" goose up -dir {{migration-dir}}; \
    fi
