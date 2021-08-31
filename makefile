build:
	go test ./...
	go build .

run:
	go run main.go

docker-build:
	docker build -t go-base .
