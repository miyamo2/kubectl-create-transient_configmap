# kubectl-create-transient_configmap 

kubectl plugin.  
Create a ConfigMap and a Job. And after the job is complete, delete them.

## Quick Start

### Install

#### With go install(recommend)

```sh
go install github.com/miyamo2/kubectl-create-transient_configmap 
```

### Simple Usage

**Manifest**
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: foo-batch
spec:
  timeZone: "Asia/Tokyo"
  schedule: "0 0 * * *"
  startingDeadlineSeconds: 100
  jobTemplate:
    spec:
      completions: 1
      parallelism: 1
      backoffLimit: 0
      template:
        spec:
          containers:
            - name: foo-batch
              image: e2e:latest
              imagePullPolicy: Never
              env:
                - name: NUM
                  valueFrom:
                    configMapKeyRef:
                      name: foo-configmap
                      key: num
                      optional: true
          restartPolicy: Never
```

**command**
```sh
kubectl create transient_configmap my-config --from-literal=num=1 --job-name=test-job --job-from=cronjob/a-cronjob
```

## Contributing

Feel free to open a PR or an Issue.  
However, you must promise to follow our [Code of Conduct](https://github.com/miyamo2/kubectl-create-transient_configmap/blob/main/CODE_OF_CONDUCT.md).

## License

**kubectl-create-transient_configmap** released under the [MIT License](https://github.com/miyamo2/kubectl-create-transient_configmap/blob/main/LICENSE)