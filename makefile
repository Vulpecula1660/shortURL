DB_URL=postgresql://root:secret@localhost:5432/short_url?sslmode=disable
makeFileDir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

network:
	docker network create shorturl-network

postgres:
	docker run --name postgres --network shorturl-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

dropdb:
	docker exec -it postgres dropdb short_url

createdb:
	docker exec -it postgres createdb --username=root --owner=root short_url

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	docker run --rm -v $(makeFileDir):/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -source ./db/sqlc/querier.go -destination ./db/mock/querier.go -package mockdb
	
.PHONY: network postgres createdb dropdb migrateup migratedown sqlc test server mock