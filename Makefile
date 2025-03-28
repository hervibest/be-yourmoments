DB_URL=postgres://postgres:postgres@localhost:5432/photo_svc?sslmode=disable&TimeZone=Asia/Jakarta
MIGRATIONS_DIR=db/migrations
PROTO_DIR=photo-svc/internal/pb
PROTO_FILE=photo.proto

.PHONY: migrate-down proto

start-photo-svc:
	cd photo-svc/cmd/web && go run main.go

start-upload-svc:
	cd upload-svc/cmd/web && go run main.go

migrate-down:
	cd photo-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-up:
	cd photo-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

proto:
	cd $(PROTO_DIR) && protoc --go_out=. --go-grpc_out=. $(PROTO_FILE)