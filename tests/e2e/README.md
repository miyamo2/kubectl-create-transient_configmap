# kubectl-create-transient_configmap /tests/e2e

e2e tests for kubectl-create-transient_configmap 

## Run tests

```sh
cd ../../ && go install . && cd tests/e2e
minikube start
eval $(minikube docker-env)
./test.sh 2 0
./test.sh 1 1
minikube stop
```