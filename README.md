# Simple Todolist deployed on KinD

## Running KinD locally with a local registry

Since I don't have a container registry online, I had to use a local registry. To enable the local registry and run KinD, simply run the following script/command:

```shell
sh kind/kind-with-registry.sh
```

## Building the Docker Container

Simply run:

```shell
docker build -t localhost:5000/todolist:latest .
docker push localhost:5000/todolist:latest
```

## Installing the Postgresql Helm Chart

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-postgresql bitnami/postgresql --version 14.3.1 --set auth.database=todo
```

## Installing the Custom HelmChart

### Pre-requisites

Install the crd for server-monitor and prometheus-rules.

```shell
kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/release-0.50/example/prometheus-operator-crd/monitoring.coreos.com_servicemonitors.yaml
kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/release-0.50/example/prometheus-operator-crd/monitoring.coreos.com_prometheusrules.yaml
```

### Install Helm Chart

```shell
helm install my-todo todolist/
```

### Exposing the port

```shell
 kubectl port-forward services/my-todo-todolist 8080:8080
```


## Notes

I couldn't run `act` locally due to some issues with their docker registry, so I couldn't run the Github Actions workflow locally.