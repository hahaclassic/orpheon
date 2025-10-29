gen-proto:
	protoc --proto_path=../../api/common/jwt/v1/proto \
		--go_out=./gen/proto --go_opt=paths=source_relative \
		--go-grpc_out=./gen/proto --go-grpc_opt=paths=source_relative \
		jwt.proto

	protoc --proto_path=../../api \
		--go_out=./gen/proto --go_opt=paths=source_relative \
		--go-grpc_out=./gen/proto --go-grpc_opt=paths=source_relative \
		../../api/auth-msv/v1/proto/auth.proto