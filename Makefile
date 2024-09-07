
run: 
	echo "Not yet complete"

build:
	go build

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

test:
	go test ./...
