# go-base
Base template for other go projects

## Requirements

**Docker (+ Compose)**
1. https://docs.docker.com/compose/install/
2. https://docs.docker.com/compose/gettingstarted/

**Kubernetes (+ Kompose)**
1. https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/
2. https://computingforgeeks.com/how-to-install-minikube-on-ubuntu-debian-linux/
3. https://kubernetes.io/docs/tasks/configure-pod-container/translate-compose-kubernetes/

**Rabbit**
1. https://x-team.com/blog/set-up-rabbitmq-with-docker-compose/

## Running

```
docker-compose up --build
```

**Sample Curl**
```
curl -X GET http://localhost:8080/v1/user/list
curl -X POST -F 'username=linuxize' -F 'email=linuxize@example.com' -F 'password=dummy' http://localhost:8080/v2/user/register
```