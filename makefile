
simplebank_container:
	docker run --name simplebank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine 

migrations:
	migrate create -ext sql -dir db/migrations -seq init_schema

createdb:
	docker exec -it simplebank createdb --username=root --owner=root bank

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate
  
test:
	go test -v -cover ./...

gitpush:
	git add . && git commit -m "sqlc" && git push

psql:
	docker exec -it simplebank psql bank

run:
	go run ./cmd/main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/flexGURU/simplebank/db/sqlc Store

.PHONY: simplebank_container migrations migrateup migratedown sqcl run mock

