include ./app.env

postgres:
	docker run --name postgres -p 5432:5432 --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18.3-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb --username=root --owner=root simple_bank

migrateup:
	migrate -path db/migration -database $(DB_SOURCE) -verbose up

migrateup1:
	migrate -path db/migration -database $(DB_SOURCE) -verbose up 1

migratedown:
	migrate -path db/migration -database $(DB_SOURCE) -verbose down

migratedown1:
	migrate -path db/migration -database $(DB_SOURCE) -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

test:
	go test -v -cover ./...

sqlc:
	sqlc generate

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/cristianemek/go-simplebank/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml


.PHONY: createdb, dropdb, postgres, migrateup, migratedown, test, server, mock, migrateup1, migratedown1, new_migration, sqlc