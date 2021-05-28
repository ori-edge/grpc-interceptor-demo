# Run protoc to generate our Go code based on the contents of
# ./api/edge_location.proto
protoc:
	protoc --go_out="pkg" --go-grpc_out="pkg" \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		api/edge_location.proto

# Build the client and server binaries - they can then be found in the bin/
# directory
build:
	go build -o bin/edge-client ./cmd/client
	go build -o bin/edge-server ./cmd/server