build:
	go mod tidy
	go test ./... -v -coverpkg=./...
	go build .

test:
	go test ./... -v -coverpkg=./...

run:
	go run main.go

d-build:
	docker build -t go-base .

c-up:
	docker-compose up --build

c-down:
	docker-compose down

c-reset-vol:
	docker-compose down --volumes

build-mocks:
	go get github.com/golang/mock/gomock
	go mod tidy
	~/go/bin/mockgen -source=usecase/user/interface.go -destination=usecase/user/mock/service.go -package=mock