postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18.3-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb --username=root --owner=root simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/cristianemek/go-simplebank/db/sqlc Store

.PHONY: createdb, dropdb, postgres, migrateup, migratedown, test, server, mock