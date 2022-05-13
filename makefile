DB_URL=postgresql://root:secret@localhost:5432/short_url?sslmode=disable

network:
	docker network create shorturl-network

postgres:
	docker run --name postgres --network shorturl-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

dropdb:
	docker exec -it postgres dropdb short_url

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root short-url

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
	
.PHONY: network postgres createdb dropdb migrateup migratedown