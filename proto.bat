
protoc -I proto -I third_party/grpc-gateway/ -I third_party/googleapis --go_out=plugins=grpc,paths=source_relative:proto/users_pb --validate_out=lang=go:proto/users_pb --grpc-gateway_out=paths=source_relative:proto/users_pb --openapiv2_out=third_party/OpenAPI/ proto/users.proto
protoc -I proto -I third_party/grpc-gateway/ -I third_party/googleapis --go_out=plugins=grpc,paths=source_relative:proto/feedback_pb --validate_out=lang=go:proto/feedback_pb --grpc-gateway_out=paths=source_relative:proto/feedback_pb --openapiv2_out=third_party/OpenAPI/ proto/feedback.proto
protoc -I proto -I third_party/grpc-gateway/ -I third_party/googleapis --go_out=plugins=grpc,paths=source_relative:proto  --grpc-gateway_out=paths=source_relative:proto  proto/enum.proto


protoc -I proto -I third_party/grpc-gateway/ -I third_party/googleapis --go_out=plugins=grpc,paths=source_relative:proto/cameras_pb --validate_out=lang=go:proto/cameras_pb --grpc-gateway_out=paths=source_relative:proto/cameras_pb --openapiv2_out=third_party/OpenAPI/ proto/cameras.proto

protoc -I proto -I third_party/grpc-gateway/ -I third_party/googleapis --go_out=plugins=grpc,paths=source_relative:proto/camera_messages_pb --validate_out=lang=go:proto/camera_messages_pb --grpc-gateway_out=paths=source_relative:proto/camera_messages_pb --openapiv2_out=third_party/OpenAPI/ proto/camera_messages.proto







statik -m -f -src third_party/OpenAPI/