install:
	go get \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		github.com/envoyproxy/protoc-gen-validate \
		github.com/rakyll/statik

gen-proto:
	protoc \
		-I proto \
		-I third_party/grpc-gateway/ \
		-I third_party/googleapis \
		--go_out=plugins=grpc,paths=source_relative:proto/users_pb \
		--validate_out=lang=go:proto/users_pb \
		--grpc-gateway_out=paths=source_relative:proto/users_pb \
		--openapiv2_out=third_party/OpenAPI/ \
		proto/users.proto

statik:
	statik -m -f -src third_party/OpenAPI/