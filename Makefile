DB_URL=postgres://postgres:postgres@localhost:5432/photo_svc?sslmode=disable&TimeZone=Asia/Jakarta
MIGRATIONS_DIR=db/migrations
PROTO_DIR=upload-svc/internal/pb
PROTO_FILE=photo.proto

.PHONY: migrate-down proto

migrate-down:
	cd upload-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

proto:
	cd $(PROTO_DIR) && protoc --go_out=. --go-grpc_out=. $(PROTO_FILE)