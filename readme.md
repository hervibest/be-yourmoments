protoc --go_out=. --go-grpc_out=. shipment.proto

goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5432/photo_svc?sslmode=disable&TimeZone=Asia/Jakarta" up\

goose -dir db/migrations create create_photo_table sql