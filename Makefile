clean: 
	go clean -testcache ...
	
build:
	go build ./...

test:
	go test -v ./...

run:
	go run main.go
