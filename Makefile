
build:
	go build -o bin/app rest/cmd/app/main.go

run: build
	./bin/app