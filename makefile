
simplebank_container:
	docker run --name simplebank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine 

migrations:
	migrate create -ext sql -dir db/migrations -seq init_schema

createdb:
	docker exec -it simplebank createdb --username=root --owner=root bank

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down 1

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

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

swagger:
	swag init -g ./cmd/main.go -o docs

evans: 
	evans --host localhost --port 9090 -r repl

.PHONY: simplebank_container migrations migrateup migratedown sqcl run mock migrateup1 migratedown1 proto evans swagger

