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

**PgBouncer**
1. https://www.compose.com/articles/how-to-pool-postgresql-connections-with-pgbouncer/
2. https://hub.docker.com/r/edoburu/pgbouncer
3. `psql postgres://postgres-dev:password@localhost/pgbouncer`

**Vegeta**
1. https://www.scaleway.com/en/docs/tutorials/load-testing-vegeta/

## Running

```
docker-compose up --build
```

**Using K8s**
1. Update run `kompose convert` to update the k8s files
2. Update files and also the `kube-apply` section in the `Makefile` with the latest k8s files
3. run `make kube-start` (if it's not running yet)
4. Login to your docker account
5. run `make kube-apply`
6. run `make kube-list` to list out the service running
7. use the returned `URL` accordingly

**Sample Curl**
```
curl -X GET http://localhost:8080/v1/user/list
curl -X POST -F 'username=linuxize' -F 'email=linuxize@example.com' -F 'password=dummy' http://localhost:8080/v2/user/register

```