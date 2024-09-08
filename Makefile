DB_URL=postgresql://admin:password@localhost:5432/invoice-go?sslmode=disable

createdb:
	docker exec -it postgresql createdb --username=admin --owner=admin invoice-go

dropdb:
	docker exec -it postgresql dropdb invoice-go

migrateup:
	migrate -path migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path migration -database "$(DB_URL)" -verbose down


new_migration:
	migrate create -ext sql -dir migration -seq $(name)

test:
	go test -v -cover -short ./...

server:
	go run cmd/api/main.go

up:
	docker-compose up -d

down:
	docker-compose down


# redis:
# 	docker run --name redis -p 6379:6379 -d redis:7-alpine

# mock:
# 	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store
# 	mockgen -package mockwk -destination worker/mock/distributor.go github.com/techschool/simplebank/worker TaskDistributor



.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration test server mock