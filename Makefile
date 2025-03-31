PHOTO_DB_URL=postgres://postgres:postgres@localhost:5432/photo_svc?sslmode=disable&TimeZone=Asia/Jakarta
USER_DB_URL=postgres://postgres:postgres@localhost:5432/user_svc?sslmode=disable&TimeZone=Asia/Jakarta
MIGRATIONS_DIR=db/migrations
PROTO_DIR=photo-svc/internal/pb
PROTO_FILE=photo.proto

.PHONY: migrate-down proto

start-photo-svc:
	cd photo-svc/cmd/web && go run main.go

start-upload-svc:
	cd upload-svc/cmd/web && go run main.go

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

photo-svc-migrate-down:
	cd photo-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(PHOTO_DB_URL)" down

photo-svc-migrate-up:
	cd photo-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(PHOTO_DB_URL)" up

user-svc-migrate-down:
	cd user-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" down

user-svc-migrate-up:
	cd user-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" up

proto:
	cd $(PROTO_DIR) && protoc --go_out=. --go-grpc_out=. $(PROTO_FILE)

mockgen:
	mockgen -source=./repository/reset_password_repository.go -destination=./mocks/repository/mock_reset_password_repository.go -package=mockrepository