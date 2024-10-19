run: 
	echo "Not yet complete"

run-docker: 
	sudo docker run -v $(realpath ./assets/latex/output):/app/assets/latex/output anemon:latest

clean:
	sudo rm ./assets/latex/output/*

build:
	go build

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

tidy:
	go mod tidy

test:
	go test -v ./...
