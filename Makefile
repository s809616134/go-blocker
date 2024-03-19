build:
	@go build -o bin/blocker

run: build
	@./bin/blocker

test: 
	@go test -v ./...

proto:
	# do this in bash before make proto
	# export PATH="$PATH:$(go env GOPATH)/bin"
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	proto/*.proto

.PHONY: proto