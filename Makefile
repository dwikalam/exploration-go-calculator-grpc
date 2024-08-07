GOPATH = $(shell go env GOPATH)

bin_dir = bin/calcgorpc

api_entry_dir = cmd/api
api_bin_dir = $(bin_dir)/server

client_entry_dir = cmd/client
client_bin_dir = $(bin_dir)/client

pb_parent_dir = internal/app
proto_dir = $(pb_parent_dir)/proto

generate-protobuf:
	@protoc --proto_path=$(proto_dir) $(proto_dir)/*.proto \
		--go_out=$(pb_parent_dir) --plugin=protoc-gen-go=$(GOPATH)/bin/protoc-gen-go \
		--go-grpc_out=$(pb_parent_dir) --plugin=protoc-gen-go-grpc=$(GOPATH)/bin/protoc-gen-go-grpc

build-server:
	@go build -o $(api_bin_dir) $(api_entry_dir)/main.go

run-server: build-server
	@./$(api_bin_dir)

build-client:
	@go build -o $(client_bin_dir) $(client_entry_dir)/main.go

run-client: build-client
	@./$(client_bin_dir)