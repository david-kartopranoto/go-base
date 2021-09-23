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
	docker-compose down --remove-orphans

c-reset-vol:
	docker-compose down --volumes

build-mocks:
	go get github.com/golang/mock/gomock
	go mod tidy
	~/go/bin/mockgen -source=usecase/user/interface.go -destination=usecase/user/mock/service.go -package=mock

kube-start:
	minikube start

kube-apply:
	docker build -t localhost/app -f Dockerfile .
	docker image tag localhost/app:latest dkartopr/app:latest
	docker build -t localhost/db -f db.Dockerfile .
	docker image tag localhost/db:latest dkartopr/db:latest
	docker login
	docker push dkartopr/app:latest
	docker push dkartopr/db:latest
	kubectl apply -f app-deployment.yaml
	kubectl apply -f app-service.yaml
	kubectl apply -f env-app-configmap.yaml
	kubectl apply -f env-postgres-configmap.yaml
	kubectl apply -f postgres-db-deployment.yaml
	kubectl apply -f postgres-db-service.yaml