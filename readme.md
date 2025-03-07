protoc --go_out=. --go-grpc_out=. shipment.proto

goose -dir db/migrations postgres "postgres://root:password@localhost:5433/photo-svc?sslmode=disable&TimeZone=Asia/Jakarta" up\

goose -dir db/migrations create create_photo_table sql