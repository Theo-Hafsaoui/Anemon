
run: 
	echo "Not yet complete"

build:
	go build

lint:
	golangci-lint run ./...

test:
	go test ./...
