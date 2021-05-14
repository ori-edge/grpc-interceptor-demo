protoc:
	protoc --go_out="pkg" --go-grpc_out="pkg" \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		api/edge_location.proto

build:
	go build -o bin/edge-client ./cmd/client
	go build -o bin/edge-server ./cmd/server