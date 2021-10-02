build-app:
	go mod tidy
	go test ./... -v -coverpkg=./...
	go build -o app ./cmd/http/main.go

test:
	go test ./... -v -coverpkg=./...

run-app:
	go run ./cmd/http/main.go

d-build-app:
	docker build -t localhost/app -f app.Dockerfile .

d-run-app:
	docker run --publish 8080:8080 localhost/app
	
d-prune:
	docker rm `docker ps --no-trunc -aq`
	docker image prune --all --filter "until=120h"
	docker images | grep none | awk '{ print $3; }' | xargs docker rmi

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
	docker build -t localhost/app -f app.Dockerfile .
	docker image tag localhost/app:latest dkartopr/app:latest
	docker build -t localhost/consumer -f consumer.Dockerfile .
	docker image tag localhost/consumer:latest dkartopr/consumer:latest
	docker build -t localhost/db -f db.Dockerfile .
	docker image tag localhost/db:latest dkartopr/db:latest
	docker login
	docker push dkartopr/app:latest
	docker push dkartopr/consumer:latest
	docker push dkartopr/db:latest
	kubectl apply -f app-service.yaml
	kubectl apply -f bouncer-service.yaml
	kubectl apply -f consumer-service.yaml
	kubectl apply -f postgres-db-service.yaml
	kubectl apply -f rabbitmq-service.yaml
	kubectl apply -f app-deployment.yaml
	kubectl apply -f bouncer-deployment.yaml
	kubectl apply -f consumer-deployment.yaml
	kubectl apply -f postgres-db-deployment.yaml
	kubectl apply -f rabbitmq-deployment.yaml
	kubectl apply -f env-app-configmap.yaml
	kubectl apply -f env-bouncer-configmap.yaml
	kubectl apply -f env-consumer-configmap.yaml
	kubectl apply -f env-postgres-configmap.yaml
	kubectl apply -f env-rabbit-configmap.yaml

kube-list:
	minikube service list