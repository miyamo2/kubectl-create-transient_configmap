# kubectl-create-transient_configmap 

[![CI](https://github.com/miyamo2/kubectl-create-transient_configmap/actions/workflows/ci.yaml/badge.svg)](https://github.com/miyamo2/kubectl-create-transient_configmap/actions/workflows/ci.yaml)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/miyamo2/kubectl-create-transient_configmap)](https://img.shields.io/github/v/release/miyamo2/kubectl-create-transient_configmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/miyamo2/kubectl-create-transient_configmap)](https://goreportcard.com/report/github.com/miyamo2/kubectl-create-transient_configmap)
[![GitHub License](https://img.shields.io/github/license/miyamo2/kubectl-create-transient_configmap?&color=blue)](https://img.shields.io/github/license/miyamo2/kubectl-create-transient_configmap?&color=blue)

kubectl plugin.  
Create a ConfigMap and a Job. And after the job is complete, delete them.

## Quick Start

### Install

#### With homebrew

```sh
brew install miyamo2/tap/kubectl-create-transient_configmap
```

#### With go install

```sh
go install github.com/miyamo2/kubectl-create-transient_configmap@latest
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
kubectl create transient_configmap foo-configmap --from-literal=num=1 --job-name=test-job --job-from=cronjob/foo-batch
```

## Features

### Flags

| name            | description                                                                                                                                                                                                                                                                                               |
|-----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `from-env-file` | Specify the path to a file to read lines of key=val pairs to create a configmap.                                                                                                                                                                                                                          |
| `from-file`     | Key file can be specified using its file path, in which case file basename will be used as configmap key, or optionally with a key and file path, in which case the given key will be used. Specifying a directory will iterate each named file in the directory whose basename is a valid configmap key. |
| `from-literal`  | Specify a key and literal value to insert in configmap (i.e. mykey=somevalue)                                                                                                                                                                                                                             |
| `job-name`      | Name of job to be created. required.                                                                                                                                                                                                                                                                      |
| `job-from`      | The name of the resource to create a Job from (only cronjob is supported).                                                                                                                                                                                                                                |
| `job-image`     | Image name to run.                                                                                                                                                                                                                                                                                        |

## Contributing

Feel free to open a PR or an Issue.  
However, you must promise to follow our [Code of Conduct](https://github.com/miyamo2/kubectl-create-transient_configmap/blob/main/CODE_OF_CONDUCT.md).

## License

**kubectl-create-transient_configmap** released under the [MIT License](https://github.com/miyamo2/kubectl-create-transient_configmap/blob/main/LICENSE)