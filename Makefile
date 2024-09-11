.PHONY: protos

protos:
	protoc --go_out=reader --go_opt=paths=source_relative --go-grpc_out=reader --go-grpc_opt=paths=source_relative reader.proto
