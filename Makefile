DB_URL=postgresql://postgres:postgres123@postgres:5432/snapstore?sslmode=disable

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres123 -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres snapshots

dropdb:
	docker exec -it postgres dropdb snapshots

migrateup:
	migrate -path db/pg/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/pg/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	cd cmd-etl && go test -v -cover -short ./... && cd ..

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
